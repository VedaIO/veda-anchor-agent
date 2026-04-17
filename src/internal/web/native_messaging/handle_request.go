package native_messaging

import (
	"encoding/json"
	"log"
	"veda-anchor-agent/src/internal/ipc"
)

// handleRequest dispatches the incoming request to the appropriate handler logic.
func handleRequest(req Request, engine *ipc.EngineClient) {
	log.Printf("Processing message type: %s", req.Type)

	switch req.Type {
	case "ping":
		sendResponse(map[string]string{"type": "pong"})

	case "log_url":
		// Handle URL logging
		var payload WebLogPayload
		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			log.Printf("Error unmarshalling log_url: %v", err)
			return
		}

		log.Printf("Logging URL: %s", payload.Url)
		_, err := engine.Request("LogWebEvent", payload.Url)
		if err != nil {
			log.Printf("Error logging web event IPC: %v", err)
		}

	case "log_web_metadata":
		var payload WebMetadataPayload
		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			log.Printf("Error unmarshalling log_web_metadata payload: %v", err)
			return
		}

		// Log metadata via IPC
		_, err := engine.Request("SaveWebMetadata", map[string]string{
			"domain":  payload.Domain,
			"title":   payload.Title,
			"iconURL": payload.IconURL,
		})
		if err != nil {
			log.Printf("Error saving metadata IPC: %v", err)
		}

	case "get_web_blocklist":
		// Send blocklist via IPC
		var bl []string
		resBytes, err := engine.Request("GetWebBlocklist", nil)
		if err != nil {
			log.Printf("Error loading blocklist IPC: %v", err)
		} else if resBytes != nil {
			json.Unmarshal(resBytes, &bl)
		}
		if bl == nil {
			bl = []string{}
		}
		sendResponse(map[string]interface{}{
			"type":    "web_blocklist",
			"payload": bl,
		})
	case "add_to_web_blocklist":
		var domain string
		if err := json.Unmarshal(req.Payload, &domain); err != nil {
			log.Printf("Error unmarshalling add_to_web_blocklist payload: %v", err)
			return
		}
		if _, err := engine.Request("AddWebBlocklist", domain); err != nil {
			log.Printf("Error adding to web blocklist: %v", err)
		}
	default:
		log.Printf("Unknown message type: %s", req.Type)
	}
}
