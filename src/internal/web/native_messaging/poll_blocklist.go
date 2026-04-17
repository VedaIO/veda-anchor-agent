package native_messaging

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"veda-anchor-agent/src/internal/ipc"
)

const (
	// pollInterval is the interval at which the web blocklist is polled for changes.
	pollInterval = 5 * time.Second
)

// pollWebBlocklist periodically checks for changes in the web blocklist and sends updates to the extension.
func pollWebBlocklist(engine *ipc.EngineClient) {
	var lastBlocklist []string
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		type blockedWebsite struct {
			Domain  string `json:"domain"`
			Title   string `json:"title"`
			IconURL string `json:"iconURL"`
		}
		var rawList []blockedWebsite
		resBytes, err := engine.Request("GetWebBlocklist", nil)
		if err != nil {
			log.Printf("Failed to get web blocklist via IPC: %v", err)
			continue
		}
		if resBytes != nil {
			json.Unmarshal(resBytes, &rawList)
		}

		list := make([]string, 0, len(rawList))
		for _, item := range rawList {
			list = append(list, item.Domain)
		}

		// Only send an update if the blocklist has changed.
		if !reflect.DeepEqual(list, lastBlocklist) {
			lastBlocklist = list
			sendResponse(map[string]any{
				"type":    "web_blocklist",
				"payload": list,
			})
		}
	}
}
