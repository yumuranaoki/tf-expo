package presentation

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/yumuranaoki/tfx/internal/model"
)

func SelectResource(resources []model.ResourceChange) (model.ResourceChange, error) {
	// Set up signal handling for graceful exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Channel to receive the result from fuzzyfinder
	resultChan := make(chan struct {
		idx int
		err error
	})
	
	// Run fuzzyfinder in a goroutine
	go func() {
		idx, err := fuzzyfinder.Find(
			resources,
			func(i int) string { return resources[i].Address + " [" + strings.Join(resources[i].Actions, ",") + "]" },
			fuzzyfinder.WithPromptString("Select resource > "),
		)
		resultChan <- struct {
			idx int
			err error
		}{idx, err}
	}()
	
	// Wait for either the result or a signal
	select {
	case result := <-resultChan:
		if result.err != nil {
			return model.ResourceChange{}, result.err
		}
		return resources[result.idx], nil
	case <-sigChan:
		fmt.Println("\nExiting...")
		os.Exit(0)
		return model.ResourceChange{}, nil
	}
}
