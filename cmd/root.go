package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/yumuranaoki/tf-expo/internal/filter"
	"github.com/yumuranaoki/tf-expo/internal/parser"
	"github.com/yumuranaoki/tf-expo/internal/presentation"
)

var (
	action string
	target string
)

var rootCmd = &cobra.Command{
	Use:   "tfx",
	Short: "Visualize Terraform plan diffs",
	RunE: func(cmd *cobra.Command, args []string) error {
		planData, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		resources, err := parser.ParsePlan(planData)
		if err != nil {
			return err
		}
		filtered := filter.Filter(resources, action, target)
		selection, err := presentation.SelectResource(filtered)
		if err != nil {
			return err
		}
		presentation.ShowDiff(selection)
		return nil
	},
}

func Execute() {
	rootCmd.Flags().StringVar(&action, "action", "", "Filter by action (create, update, delete, replace)")
	rootCmd.Flags().StringVar(&target, "target", "", "Filter by module/target prefix")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
