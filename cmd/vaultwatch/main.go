// Package main is the entry point for the vaultwatch daemon.
// It wires together configuration, Vault client, webhook sender,
// and the lease monitor, then runs until the process is signalled.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/monitor"
	"github.com/yourusername/vaultwatch/internal/vault"
	"github.com/yourusername/vaultwatch/internal/webhook"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfgPath := os.Getenv("VAULTWATCH_CONFIG")
	if cfgPath == "" {
		cfgPath = "vaultwatch.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		slog.Error("failed to load configuration", "path", cfgPath, "error", err)
		os.Exit(1)
	}

	slog.Info("configuration loaded",
		"vault_address", cfg.Vault.Address,
		"check_interval", cfg.Monitor.CheckInterval,
		"warn_threshold", cfg.Monitor.WarnThreshold,
	)

	vaultClient, err := vault.NewClient(cfg)
	if err != nil {
		slog.Error("failed to create Vault client", "error", err)
		os.Exit(1)
	}

	// Verify connectivity to Vault before starting the monitor loop.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if healthy, err := vaultClient.IsHealthy(ctx); err != nil || !healthy {
		slog.Error("vault health check failed", "error", err)
		os.Exit(1)
	}
	slog.Info("vault health check passed")

	retryConfig := webhook.DefaultRetryConfig()
	webhookSender := webhook.New(cfg, retryConfig)

	mon := monitor.New(cfg, vaultClient, webhookSender)

	// Listen for OS signals so we can shut down gracefully.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		slog.Info("received shutdown signal", "signal", sig)
		cancel()
	}()

	slog.Info("starting vaultwatch monitor")
	if err := mon.Run(ctx); err != nil && err != context.Canceled {
		slog.Error("monitor exited with error", "error", err)
		os.Exit(1)
	}

	slog.Info("vaultwatch stopped")
}
