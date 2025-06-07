package filter

import (
	"testing"

	"github.com/yumuranaoki/tf-expo/internal/model"
)

func TestFilter(t *testing.T) {
	testResources := []model.ResourceChange{
		{
			Address: "aws_instance.web",
			Actions: []string{"create"},
		},
		{
			Address: "aws_instance.api",
			Actions: []string{"update"},
		},
		{
			Address: "aws_s3_bucket.data",
			Actions: []string{"delete"},
		},
		{
			Address: "module.vpc.aws_vpc.main",
			Actions: []string{"no-op"},
		},
		{
			Address: "module.database.aws_rds_instance.main",
			Actions: []string{"delete", "create"},
		},
		{
			Address: "aws_iam_role.lambda",
			Actions: []string{"update"},
		},
	}

	testCases := []struct {
		name           string
		resources      []model.ResourceChange
		action         string
		target         string
		expectedCount  int
		expectedAddrs  []string
	}{
		{
			name:          "no filters",
			resources:     testResources,
			action:        "",
			target:        "",
			expectedCount: 5, // excludes no-op
			expectedAddrs: []string{
				"aws_instance.web",
				"aws_instance.api",
				"aws_s3_bucket.data",
				"module.database.aws_rds_instance.main",
				"aws_iam_role.lambda",
			},
		},
		{
			name:          "filter by create action",
			resources:     testResources,
			action:        "create",
			target:        "",
			expectedCount: 2,
			expectedAddrs: []string{
				"aws_instance.web",
				"module.database.aws_rds_instance.main",
			},
		},
		{
			name:          "filter by update action",
			resources:     testResources,
			action:        "update",
			target:        "",
			expectedCount: 2,
			expectedAddrs: []string{
				"aws_instance.api",
				"aws_iam_role.lambda",
			},
		},
		{
			name:          "filter by delete action",
			resources:     testResources,
			action:        "delete",
			target:        "",
			expectedCount: 2,
			expectedAddrs: []string{
				"aws_s3_bucket.data",
				"module.database.aws_rds_instance.main",
			},
		},
		{
			name:          "filter by target prefix",
			resources:     testResources,
			action:        "",
			target:        "aws_instance",
			expectedCount: 2,
			expectedAddrs: []string{
				"aws_instance.web",
				"aws_instance.api",
			},
		},
		{
			name:          "filter by module target prefix",
			resources:     testResources,
			action:        "",
			target:        "module.database",
			expectedCount: 1,
			expectedAddrs: []string{
				"module.database.aws_rds_instance.main",
			},
		},
		{
			name:          "filter by action and target",
			resources:     testResources,
			action:        "update",
			target:        "aws_instance",
			expectedCount: 1,
			expectedAddrs: []string{
				"aws_instance.api",
			},
		},
		{
			name:          "no matches for action filter",
			resources:     testResources,
			action:        "replace",
			target:        "",
			expectedCount: 0,
			expectedAddrs: []string{},
		},
		{
			name:          "no matches for target filter",
			resources:     testResources,
			action:        "",
			target:        "nonexistent",
			expectedCount: 0,
			expectedAddrs: []string{},
		},
		{
			name:          "empty resource list",
			resources:     []model.ResourceChange{},
			action:        "",
			target:        "",
			expectedCount: 0,
			expectedAddrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := Filter(tc.resources, tc.action, tc.target)

			if len(result) != tc.expectedCount {
				t.Errorf("expected %d resources, got %d", tc.expectedCount, len(result))
				return
			}

			actualAddrs := make([]string, len(result))
			for i, r := range result {
				actualAddrs[i] = r.Address
			}

			for i, expectedAddr := range tc.expectedAddrs {
				if i >= len(actualAddrs) {
					t.Errorf("missing expected address: %s", expectedAddr)
					continue
				}
				if actualAddrs[i] != expectedAddr {
					t.Errorf("expected address %s at index %d, got %s", expectedAddr, i, actualAddrs[i])
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "item exists",
			slice:    []string{"create", "update", "delete"},
			item:     "update",
			expected: true,
		},
		{
			name:     "item does not exist",
			slice:    []string{"create", "update", "delete"},
			item:     "replace",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			item:     "create",
			expected: false,
		},
		{
			name:     "empty item",
			slice:    []string{"create", "update", "delete"},
			item:     "",
			expected: false,
		},
		{
			name:     "single item match",
			slice:    []string{"create"},
			item:     "create",
			expected: true,
		},
		{
			name:     "case sensitive",
			slice:    []string{"Create", "Update", "Delete"},
			item:     "create",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := contains(tc.slice, tc.item)

			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}