package main

import (
	"os"

	"github.com/magomedcoder/llm-runner/config"
	"github.com/magomedcoder/llm-runner/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.E("config: %v", err)
		os.Exit(1)
	}
	logger.Default.SetLevel(logger.ParseLevel(cfg.LogLevel))
	logger.I("started")
}
