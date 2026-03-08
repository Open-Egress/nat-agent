package main

import (
	"log/slog"
	"os"

	"nat-agent/internal/kernel"
	"nat-agent/internal/routing"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := kernel.IPForward(); err != nil {
		slog.Error("Critical error setting up kernel parameters", "error", err)
		os.Exit(1)
	}

	if err := routing.NftableApply(); err != nil {
		slog.Error("Critical error applying nftables rules", "error", err)
		os.Exit(1)
	}

	router := gin.Default()

	if err := router.Run(":7090"); err != nil {
		slog.Error("Failed to start the server", "error", err)
		os.Exit(1)
	}
}
