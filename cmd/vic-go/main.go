package main

import (
	"os"

	"github.com/vic-sdd/vic/internal/commands"
	"github.com/vic-sdd/vic/internal/config"
)

func main() {
	// Load configuration (environment variables, config file)
	cfg := config.Load()

	// Execute root command
	if err := commands.NewRootCmd(cfg).Execute(); err != nil {
		os.Exit(1)
	}
}
