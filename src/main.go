//go:generate goversioninfo -64

package main

import (
	"log"
	"os"
	"path/filepath"
	"veda-anchor-agent/src/internal/ipc"
	"veda-anchor-agent/src/internal/tracking"
)

func main() {
	// Setup logging
	logDir := filepath.Join(os.Getenv("LocalAppData"), "VedaAnchor", "logs")
	_ = os.MkdirAll(logDir, 0755)
	
	logPath := filepath.Join(logDir, "veda-anchor-agent.log")
	logFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if logFile != nil {
		defer func() { _ = logFile.Close() }()
		log.SetOutput(logFile)
	}

	log.Printf("=== VEDA ANCHOR AGENT STARTING ===")

	// Connect to Engine IPC
	engineClient, err := ipc.NewEngineClient()
	if err != nil {
		log.Printf("Failed to connect to Engine: %v", err)
		log.Printf("Agent requires Engine to be running. Exiting.")
		os.Exit(1)
	}
	log.Printf("Connected to Engine IPC")

	// Start UI IPC Server (forwards to Engine)
	uiServer := ipc.NewUIServer(engineClient)
	go func() {
		if err := uiServer.Start(); err != nil {
			log.Printf("UI IPC server error: %v", err)
		}
	}()

	// Start tracking (window + screentime)
	tracking.Start(engineClient)

	// Wait for exit signal
	log.Printf("=== VEDA ANCHOR AGENT RUNNING ===")
	select {}
}
