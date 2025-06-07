package presentation

import (
	"encoding/json"
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/yumuranaoki/tfx/internal/model"
)

func ShowDiff(r model.ResourceChange) {
	fmt.Printf("Resource: %s\n", r.Address)
	fmt.Printf("Actions: %v\n", r.Actions)

	// Handle null values appropriately
	var beforeStr, afterStr string

	if r.Before == nil {
		beforeStr = "{}" // Empty object for create operations
	} else {
		before, _ := json.MarshalIndent(r.Before, "", "  ")
		beforeStr = string(before)
	}

	if r.After == nil {
		afterStr = "{}" // Empty object for delete operations
	} else {
		after, _ := json.MarshalIndent(r.After, "", "  ")
		afterStr = string(after)
	}

	// Show colorized diff only
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(beforeStr, afterStr, false)

	// Clean up diff output for better readability
	dmp.DiffCleanupSemantic(diffs)

	fmt.Println("Diff:")
	fmt.Println(dmp.DiffPrettyText(diffs))
}
