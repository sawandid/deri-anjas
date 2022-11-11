package main

import (
	"os"

	"github.com/sawandid/deri-anjas/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
