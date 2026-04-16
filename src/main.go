//go:generate goversioninfo -64

package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"veda-anchor-agent/src/api"
	"veda-anchor-agent/src/internal/config"
	"veda-anchor-agent/src/internal/ipc"
	"veda-anchor-agent/src/internal/tracking"
)

func main() {
	// Setup logging (use ProgramData for shared access)
	logPath, err := config.GetLogPath()
	if err != nil {
		// Fallback for safety
		logPath = filepath.Join("C:\\ProgramData", "VedaAnchor", "logs", "veda-anchor-agent.log")
	}
	_ = os.MkdirAll(filepath.Dir(logPath), 0755)

	logFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if logFile != nil {
		defer func() { _ = logFile.Close() }()
		log.SetOutput(logFile)
	}

	log.Printf("=== VEDA ANCHOR AGENT STARTING ===")

	// Connect to Engine IPC (with retry)
	var engineClient *ipc.EngineClient
	for i := range 30 {
		engineClient, err = ipc.NewEngineClient()
		if err == nil {
			break
		}
		log.Printf("Failed to connect to Engine (attempt %d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if engineClient == nil {
		log.Printf("Failed to connect to Engine after 30 attempts. Exiting.")
		os.Exit(1)
	}
	log.Printf("Connected to Engine IPC")

	// Initialize API server (handles local + forwarded methods)
	apiServer := api.NewServer(engineClient)

	// Start UI IPC Server (dispatches through apiServer)
	uiServer := ipc.NewUIServer(apiServer)
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
