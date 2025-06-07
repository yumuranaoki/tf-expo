package model

import (
	"reflect"
	"testing"
)

func TestResourceChange(t *testing.T) {
	testCases := []struct {
		name     string
		resource ResourceChange
		validate func(t *testing.T, r ResourceChange)
	}{
		{
			name: "create resource",
			resource: ResourceChange{
				Address: "aws_instance.web",
				Actions: []string{"create"},
				Before:  nil,
				After: map[string]interface{}{
					"instance_type": "t2.micro",
					"ami":           "ami-12345",
				},
			},
			validate: func(t *testing.T, r ResourceChange) {
				if r.Address != "aws_instance.web" {
					t.Errorf("expected address 'aws_instance.web', got %q", r.Address)
				}
				if len(r.Actions) != 1 || r.Actions[0] != "create" {
					t.Errorf("expected actions ['create'], got %v", r.Actions)
				}
				if r.Before != nil {
					t.Errorf("expected Before to be nil, got %v", r.Before)
				}
				if r.After == nil {
					t.Error("expected After to be non-nil")
				}
			},
		},
		{
			name: "update resource",
			resource: ResourceChange{
				Address: "aws_instance.api",
				Actions: []string{"update"},
				Before: map[string]interface{}{
					"instance_type": "t2.micro",
				},
				After: map[string]interface{}{
					"instance_type": "t2.small",
				},
			},
			validate: func(t *testing.T, r ResourceChange) {
				if r.Address != "aws_instance.api" {
					t.Errorf("expected address 'aws_instance.api', got %q", r.Address)
				}
				if len(r.Actions) != 1 || r.Actions[0] != "update" {
					t.Errorf("expected actions ['update'], got %v", r.Actions)
				}
				if r.Before == nil || r.After == nil {
					t.Error("expected both Before and After to be non-nil for update")
				}
			},
		},
		{
			name: "delete resource",
			resource: ResourceChange{
				Address: "aws_s3_bucket.old",
				Actions: []string{"delete"},
				Before: map[string]interface{}{
					"bucket": "my-old-bucket",
				},
				After: nil,
			},
			validate: func(t *testing.T, r ResourceChange) {
				if r.Address != "aws_s3_bucket.old" {
					t.Errorf("expected address 'aws_s3_bucket.old', got %q", r.Address)
				}
				if len(r.Actions) != 1 || r.Actions[0] != "delete" {
					t.Errorf("expected actions ['delete'], got %v", r.Actions)
				}
				if r.Before == nil {
					t.Error("expected Before to be non-nil for delete")
				}
				if r.After != nil {
					t.Errorf("expected After to be nil, got %v", r.After)
				}
			},
		},
		{
			name: "replace resource (delete + create)",
			resource: ResourceChange{
				Address: "aws_instance.replace_me",
				Actions: []string{"delete", "create"},
				Before: map[string]interface{}{
					"ami": "ami-old",
				},
				After: map[string]interface{}{
					"ami": "ami-new",
				},
			},
			validate: func(t *testing.T, r ResourceChange) {
				if r.Address != "aws_instance.replace_me" {
					t.Errorf("expected address 'aws_instance.replace_me', got %q", r.Address)
				}
				expectedActions := []string{"delete", "create"}
				if !reflect.DeepEqual(r.Actions, expectedActions) {
					t.Errorf("expected actions %v, got %v", expectedActions, r.Actions)
				}
				if r.Before == nil || r.After == nil {
					t.Error("expected both Before and After to be non-nil for replace")
				}
			},
		},
		{
			name: "no-op resource",
			resource: ResourceChange{
				Address: "aws_vpc.main",
				Actions: []string{"no-op"},
				Before: map[string]interface{}{
					"cidr_block": "10.0.0.0/16",
				},
				After: map[string]interface{}{
					"cidr_block": "10.0.0.0/16",
				},
			},
			validate: func(t *testing.T, r ResourceChange) {
				if r.Address != "aws_vpc.main" {
					t.Errorf("expected address 'aws_vpc.main', got %q", r.Address)
				}
				if len(r.Actions) != 1 || r.Actions[0] != "no-op" {
					t.Errorf("expected actions ['no-op'], got %v", r.Actions)
				}
				if !reflect.DeepEqual(r.Before, r.After) {
					t.Error("expected Before and After to be equal for no-op")
				}
			},
		},
		{
			name: "module resource",
			resource: ResourceChange{
				Address: "module.vpc.aws_subnet.private[0]",
				Actions: []string{"create"},
				Before:  nil,
				After: map[string]interface{}{
					"cidr_block": "10.0.1.0/24",
				},
			},
			validate: func(t *testing.T, r ResourceChange) {
				if r.Address != "module.vpc.aws_subnet.private[0]" {
					t.Errorf("expected address 'module.vpc.aws_subnet.private[0]', got %q", r.Address)
				}
			},
		},
		{
			name: "empty resource change",
			resource: ResourceChange{
				Address: "",
				Actions: []string{},
				Before:  nil,
				After:   nil,
			},
			validate: func(t *testing.T, r ResourceChange) {
				if r.Address != "" {
					t.Errorf("expected empty address, got %q", r.Address)
				}
				if len(r.Actions) != 0 {
					t.Errorf("expected empty actions, got %v", r.Actions)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.validate(t, tc.resource)
		})
	}
}