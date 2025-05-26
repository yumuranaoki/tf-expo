package presentation

import (
	"encoding/json"
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/yumuranaoki/tf-expo/internal/model"
)

func ShowDiff(r model.ResourceChange) {
	before, _ := json.MarshalIndent(r.Before, "", "  ")
	after, _ := json.MarshalIndent(r.After, "", "  ")
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(before), string(after), false)
	fmt.Println(dmp.DiffPrettyText(diffs))
}
