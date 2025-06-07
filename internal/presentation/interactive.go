package presentation

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yumuranaoki/tfx/internal/model"
)

func RunInteractiveMode(resources []model.ResourceChange) error {
	// Set up signal handling for graceful exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		selection, err := SelectResource(resources)
		if err != nil {
			return err
		}

		// Clear the screen and show diff immediately
		fmt.Print("\033[2J\033[H")
		fmt.Printf("Selected resource: %s\n\n", selection.Address)
		ShowDiff(selection)

		fmt.Println("\n\nPress Enter to go back to resource list, or 'q' + Enter to quit")
		
		// Create a channel for user input
		inputChan := make(chan bool)
		
		go func() {
			// Wait for any input from /dev/tty to avoid stdin conflicts
			tty, err := os.OpenFile("/dev/tty", os.O_RDONLY, 0)
			if err != nil {
				// Fallback: just wait for any key from stdin
				var dummy [1]byte
				os.Stdin.Read(dummy[:])
			} else {
				var buf [10]byte
				n, _ := tty.Read(buf[:])
				tty.Close()
				
				// Check if user typed specific keys
				if n > 0 {
					if buf[0] == 'q' || buf[0] == 'Q' {
						inputChan <- true // quit
						return
					} else if buf[0] == 10 || buf[0] == 13 { // Enter key (LF or CR)
						inputChan <- false // go back to list
						return
					}
				}
			}
			inputChan <- false // continue (any other key goes back)
		}()

		// Wait for either user input or signal
		select {
		case shouldQuit := <-inputChan:
			if shouldQuit {
				return nil
			}
		case <-sigChan:
			fmt.Println("\nExiting...")
			return nil
		}
	}
}