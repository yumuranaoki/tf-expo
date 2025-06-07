package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yumuranaoki/tf-expo/internal/parser"
)

func TestExecute(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		stdin       string
		expectError bool
		skipRun     bool
	}{
		{
			name:    "help flag",
			args:    []string{"--help"},
			stdin:   "",
			skipRun: true,
		},
		{
			name:        "invalid json input",
			args:        []string{},
			stdin:       `{invalid json}`,
			expectError: true,
		},
		{
			name:    "valid empty plan",
			args:    []string{},
			stdin:   `{"resource_changes": []}`,
			skipRun: true,
		},
		{
			name: "valid plan with action filter",
			args: []string{"--action", "create"},
			stdin: `{
				"resource_changes": [
					{
						"address": "aws_instance.web",
						"change": {
							"actions": ["create"],
							"before": null,
							"after": {"instance_type": "t2.micro"}
						}
					}
				]
			}`,
			skipRun: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skipRun {
				t.Skip("Skipping interactive test")
			}

			oldStdin := os.Stdin
			defer func() {
				os.Stdin = oldStdin
			}()

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}
			os.Stdin = r

			go func() {
				defer w.Close()
				io.WriteString(w, tc.stdin)
			}()

			cmd := &cobra.Command{
				Use:   "tfx",
				Short: "Visualize Terraform plan diffs",
				RunE: func(cmd *cobra.Command, args []string) error {
					planData, err := io.ReadAll(os.Stdin)
					if err != nil {
						return err
					}
					_, err = parser.ParsePlan(planData)
					return err
				},
			}
			cmd.Flags().StringVar(&action, "action", "", "Filter by action (create, update, delete, replace)")
			cmd.Flags().StringVar(&target, "target", "", "Filter by module/target prefix")
			cmd.SetArgs(tc.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			err = cmd.Execute()

			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRootCmdFlags(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedAction string
		expectedTarget string
	}{
		{
			name:           "no flags",
			args:           []string{},
			expectedAction: "",
			expectedTarget: "",
		},
		{
			name:           "action flag only",
			args:           []string{"--action", "create"},
			expectedAction: "create",
			expectedTarget: "",
		},
		{
			name:           "target flag only",
			args:           []string{"--target", "aws_instance"},
			expectedAction: "",
			expectedTarget: "aws_instance",
		},
		{
			name:           "both flags",
			args:           []string{"--action", "update", "--target", "module.vpc"},
			expectedAction: "update",
			expectedTarget: "module.vpc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var testAction, testTarget string

			cmd := &cobra.Command{
				Use:   "tfx",
				Short: "Visualize Terraform plan diffs",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}
			cmd.Flags().StringVar(&testAction, "action", "", "Filter by action (create, update, delete, replace)")
			cmd.Flags().StringVar(&testTarget, "target", "", "Filter by module/target prefix")
			cmd.SetArgs(tc.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			cmd.Execute()

			if testAction != tc.expectedAction {
				t.Errorf("expected action %q, got %q", tc.expectedAction, testAction)
			}
			if testTarget != tc.expectedTarget {
				t.Errorf("expected target %q, got %q", tc.expectedTarget, testTarget)
			}
		})
	}
}

func TestRootCmdUsage(t *testing.T) {
	var buf bytes.Buffer
	var testAction, testTarget string

	cmd := &cobra.Command{
		Use:   "tfx",
		Short: "Visualize Terraform plan diffs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.Flags().StringVar(&testAction, "action", "", "Filter by action (create, update, delete, replace)")
	cmd.Flags().StringVar(&testTarget, "target", "", "Filter by module/target prefix")
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "help") {
		if err == nil {
			output := buf.String()
			if !strings.Contains(output, "Usage:") {
				t.Error("expected help output but got none")
			}
		}
	}

	output := buf.String()
	expectedStrings := []string{
		"tfx",
		"Visualize Terraform plan diffs",
		"--action",
		"--target",
		"Filter by action",
		"Filter by module/target prefix",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("expected help output to contain %q, but it didn't. Output:\n%s", expected, output)
		}
	}
}