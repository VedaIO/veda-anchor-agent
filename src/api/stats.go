package api

import (
	"encoding/json"
)

// --- Types ---

type AppLeaderboardItem struct {
	Rank           int    `json:"rank"`
	Name           string `json:"name"`           // Display name (commercial name if available)
	ProcessName    string `json:"processName"`    // Actual process name for blocking
	ExecutablePath string `json:"executablePath"` // Used locally for icon enrichment
	Icon           string `json:"icon"`
	Count          int    `json:"count"`
}

type WebLeaderboardItem struct {
	Rank   int    `json:"rank"`
	Domain string `json:"domain"`
	Title  string `json:"title"`
	Icon   string `json:"icon"`
	Count  int    `json:"count"`
}

type ScreenTimeItem struct {
	Name            string `json:"name"`
	ExecutablePath  string `json:"executablePath"`
	Icon            string `json:"icon"`
	DurationSeconds int    `json:"durationSeconds"`
}

// --- App Usage ---

func (s *Server) GetAppLeaderboard(since, until string) (json.RawMessage, error) {
	return s.engine.Request("GetAppLeaderboard", map[string]string{"since": since, "until": until})
}

func (s *Server) GetScreenTime() (json.RawMessage, error) {
	return s.engine.Request("GetScreenTime", nil)
}

func (s *Server) GetTotalScreenTime() (json.RawMessage, error) {
	return s.engine.Request("GetTotalScreenTime", nil)
}

// --- Web Usage ---

func (s *Server) GetWebLeaderboard(since, until string) (json.RawMessage, error) {
	return s.engine.Request("GetWebLeaderboard", map[string]string{"since": since, "until": until})
}

// --- Logs & Search ---

func (s *Server) Search(query, since, until string) (json.RawMessage, error) {
	return s.engine.Request("Search", map[string]string{"query": query, "since": since, "until": until})
}

func (s *Server) GetWebLogs(query, since, until string) (json.RawMessage, error) {
	return s.engine.Request("GetWebLogs", map[string]string{"query": query, "since": since, "until": until})
}
