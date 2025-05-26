package filter

import (
	"strings"

	"github.com/yumuranaoki/tf-expo/internal/model"
)

func Filter(resources []model.ResourceChange, action, target string) []model.ResourceChange {
	var result []model.ResourceChange
	for _, r := range resources {
		if contains(r.Actions, "no-op") {
			continue
		}
		if action != "" && !contains(r.Actions, action) {
			continue
		}
		if target != "" && !strings.HasPrefix(r.Address, target) {
			continue
		}
		result = append(result, r)
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
