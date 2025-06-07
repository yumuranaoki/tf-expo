package parser

import (
	"testing"

	"github.com/yumuranaoki/tf-expo/internal/model"
)

func TestParsePlan(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected []model.ResourceChange
		wantErr  bool
	}{
		{
			name: "valid plan with single resource",
			input: []byte(`{
				"resource_changes": [
					{
						"address": "aws_instance.example",
						"change": {
							"actions": ["create"],
							"before": null,
							"after": {
								"instance_type": "t2.micro"
							}
						}
					}
				]
			}`),
			expected: []model.ResourceChange{
				{
					Address: "aws_instance.example",
					Actions: []string{"create"},
					Before:  nil,
					After: map[string]interface{}{
						"instance_type": "t2.micro",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid plan with multiple resources",
			input: []byte(`{
				"resource_changes": [
					{
						"address": "aws_instance.web",
						"change": {
							"actions": ["update"],
							"before": {
								"instance_type": "t2.micro"
							},
							"after": {
								"instance_type": "t2.small"
							}
						}
					},
					{
						"address": "aws_s3_bucket.data",
						"change": {
							"actions": ["delete"],
							"before": {
								"bucket": "my-bucket"
							},
							"after": null
						}
					}
				]
			}`),
			expected: []model.ResourceChange{
				{
					Address: "aws_instance.web",
					Actions: []string{"update"},
					Before: map[string]interface{}{
						"instance_type": "t2.micro",
					},
					After: map[string]interface{}{
						"instance_type": "t2.small",
					},
				},
				{
					Address: "aws_s3_bucket.data",
					Actions: []string{"delete"},
					Before: map[string]interface{}{
						"bucket": "my-bucket",
					},
					After: nil,
				},
			},
			wantErr: false,
		},
		{
			name:     "empty plan",
			input:    []byte(`{"resource_changes": []}`),
			expected: []model.ResourceChange{},
			wantErr:  false,
		},
		{
			name:     "invalid JSON",
			input:    []byte(`{invalid json}`),
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "empty input",
			input:    []byte(``),
			expected: nil,
			wantErr:  true,
		},
		{
			name: "resource with multiple actions",
			input: []byte(`{
				"resource_changes": [
					{
						"address": "aws_instance.replace_me",
						"change": {
							"actions": ["delete", "create"],
							"before": {
								"ami": "ami-123"
							},
							"after": {
								"ami": "ami-456"
							}
						}
					}
				]
			}`),
			expected: []model.ResourceChange{
				{
					Address: "aws_instance.replace_me",
					Actions: []string{"delete", "create"},
					Before: map[string]interface{}{
						"ami": "ami-123",
					},
					After: map[string]interface{}{
						"ami": "ami-456",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := ParsePlan(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != len(tc.expected) {
				t.Errorf("expected %d resources, got %d", len(tc.expected), len(result))
				return
			}

			for i, expected := range tc.expected {
				actual := result[i]
				if actual.Address != expected.Address {
					t.Errorf("resource %d: expected address %q, got %q", i, expected.Address, actual.Address)
				}

				if len(actual.Actions) != len(expected.Actions) {
					t.Errorf("resource %d: expected %d actions, got %d", i, len(expected.Actions), len(actual.Actions))
					continue
				}

				for j, expectedAction := range expected.Actions {
					if actual.Actions[j] != expectedAction {
						t.Errorf("resource %d action %d: expected %q, got %q", i, j, expectedAction, actual.Actions[j])
					}
				}
			}
		})
	}
}