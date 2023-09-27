package main

import (
	"github.com/pygrum/siphon/cmd"
	"github.com/pygrum/siphon/internal/logger"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
