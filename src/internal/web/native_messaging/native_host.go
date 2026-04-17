package native_messaging

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
	"veda-anchor-agent/src/internal/config"
	"veda-anchor-agent/src/internal/ipc"
)

// Run starts the native messaging host loop. It sets up logging, initializes the database,
// and begins listening for messages from the browser extension via standard input.
func Run() {
	// Setup logging to file (CRITICAL for debugging native messaging)
	logPath, err := config.GetNativeHostLogPath()
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(logPath), 0755)

	logFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if logFile != nil {
		defer func() { _ = logFile.Close() }()
		log.SetOutput(logFile)
	}

	// Connect to Engine
	engine, err := ipc.NewEngineClient()
	if err != nil {
		log.Printf("CRITICAL: Failed to connect to Engine IPC: %v", err)
	} else {
		log.Println("Engine IPC connected successfully")
		defer engine.Close()
	}

	// Catch panics to see why it crashes
	defer func() {
		if r := recover(); r != nil {
			log.Printf("CRITICAL PANIC in Native Host: %v\nStack: %s", r, string(debug.Stack()))
		}
	}()

	log.Println("=== NATIVE MESSAGING HOST STARTED ===")

	// Start blocklist poller (panic safe)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC in Blocklist Poller: %v", r)
			}
		}()
		pollWebBlocklist(engine)
	}()

	// Start continuous heartbeat updater
	// This ensures the GUI knows we are alive even if no messages are flowing
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			updateHeartbeat()
		}
	}()

	// Main Message Loop
	for {
		log.Println("Waiting for message...")

		// 1. Read Message Length (4 bytes, little endian)
		var length uint32
		if err := binary.Read(os.Stdin, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				log.Println("Chrome disconnected (EOF)")
				return
			}
			log.Printf("Error reading length: %v", err)
			return
		}

		// 2. Read Message Body
		msg := make([]byte, length)
		if _, err := io.ReadFull(os.Stdin, msg); err != nil {
			log.Printf("Error reading body: %v", err)
			return
		}

		log.Printf("Received message (%d bytes): %s", length, string(msg))

		// 3. Update Heartbeat File (for GUI detection)
		updateHeartbeat()

		// 4. Process Message
		var req Request
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Printf("JSON Error: %v", err)
			continue
		}

		handleRequest(req, engine)

		log.Println("Message processed successfully")
	}
}

// Stop sends a stopping message to the extension to prevent it from reconnecting.
func Stop() {
	sendResponse(map[string]interface{}{
		"type":    "stopping",
		"payload": nil,
	})
}
