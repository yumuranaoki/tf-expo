package parser

import (
	"encoding/json"

	"github.com/yumuranaoki/tf-expo/internal/model"
)

func ParsePlan(data []byte) ([]model.ResourceChange, error) {
	var out struct {
		ResourceChanges []struct {
			Address string `json:"address"`
			Change  struct {
				Actions []string               `json:"actions"`
				Before  map[string]interface{} `json:"before"`
				After   map[string]interface{} `json:"after"`
			} `json:"change"`
		} `json:"resource_changes"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	resources := make([]model.ResourceChange, len(out.ResourceChanges))
	for i, rc := range out.ResourceChanges {
		resources[i] = model.ResourceChange{
			Address: rc.Address,
			Actions: rc.Change.Actions,
			Before:  rc.Change.Before,
			After:   rc.Change.After,
		}
	}
	return resources, nil
}
