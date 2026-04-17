//go:build windows

package tracking

import (
	"log"
	"sync"
	"time"

	"veda-anchor-agent/src/internal/ipc"
	"veda-anchor-agent/src/internal/platform/screentime"
	"veda-anchor-agent/src/internal/platform/window"

	"golang.org/x/sys/windows"
)

var (
	lastPID     uint32
	lastCheck   time.Time
	pendingSecs int64
	stateMu     sync.Mutex
)

// Reset clears the internal tracking state (used when history is cleared).
func Reset() {
	stateMu.Lock()
	defer stateMu.Unlock()
	lastPID = 0
	pendingSecs = 0
}

func Start(engineClient *ipc.EngineClient) {
	log.Println("[Tracking] Starting window and screentime tracking")

	// Start window tracking (polls for visible foreground window)
	go trackForegroundWindow(engineClient)

	// Start periodic flush to Engine
	go periodicFlush(engineClient)
}

func trackForegroundWindow(client *ipc.EngineClient) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stateMu.Lock()
		pid, err := screentime.GetForegroundPID()
		if err != nil {
			stateMu.Unlock()
			continue
		}

		// Check if process has visible window
		if pid != 0 && window.HasVisibleWindow(pid) {
			if pid != lastPID {
				// Window changed - flush previous and start new
				if lastPID != 0 && pendingSecs > 0 {
					sendScreenTimeUpdate(client, lastPID, pendingSecs)
				}
				lastPID = pid
				pendingSecs = 0

				// Report new active app to engine (creates initial screen_time record)
				if exePath, err := getExePathByPID(pid); err == nil && exePath != "" {
					reportActiveApp(client, pid, exePath)
				}
			}
			pendingSecs++
		}
		lastCheck = time.Now()
		stateMu.Unlock()
	}
}

func getExePathByPID(pid uint32) (string, error) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(h)

	var pathBuf [windows.MAX_PATH]uint16
	pathLen := uint32(len(pathBuf))
	err = windows.QueryFullProcessImageName(h, 0, &pathBuf[0], &pathLen)
	if err != nil {
		return "", err
	}
	return windows.UTF16ToString(pathBuf[:pathLen]), nil
}

func reportActiveApp(client *ipc.EngineClient, pid uint32, exePath string) {
	log.Printf("[Tracking] Reporting active app: PID=%d, ExePath=%s", pid, exePath)
	params := map[string]any{
		"pid":     pid,
		"exePath": exePath,
	}
	_, err := client.Request("ReportActiveApp", params)
	if err != nil {
		log.Printf("[Tracking] Failed to report active app: %v", err)
	}
}

func periodicFlush(client *ipc.EngineClient) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stateMu.Lock()
		if lastPID != 0 && pendingSecs > 0 {
			sendScreenTimeUpdate(client, lastPID, pendingSecs)
			pendingSecs = 0
		}
		stateMu.Unlock()
	}
}

func sendScreenTimeUpdate(client *ipc.EngineClient, pid uint32, seconds int64) {
	log.Printf("[Tracking] Sending screentime: PID=%d, Seconds=%d", pid, seconds)

	params := map[string]any{
		"pid":     pid,
		"seconds": seconds,
	}

	_, err := client.Request("UpdateScreenTime", params)
	if err != nil {
		log.Printf("[Tracking] Failed to send screentime: %v", err)
	}
}
