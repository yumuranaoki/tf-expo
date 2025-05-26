package presentation

import (
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/yumuranaoki/tf-expo/internal/model"
)

func SelectResource(resources []model.ResourceChange) (model.ResourceChange, error) {
	idx, err := fuzzyfinder.Find(
		resources,
		func(i int) string { return resources[i].Address + " [" + strings.Join(resources[i].Actions, ",") + "]" },
		fuzzyfinder.WithPromptString("Select resource > "),
	)
	if err != nil {
		return model.ResourceChange{}, err
	}
	return resources[idx], nil
}
