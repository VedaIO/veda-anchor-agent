package api

import (
	"encoding/json"
	"log"

	"veda-anchor-agent/src/internal/platform/nativehost"
)

// --- Lifecycle ---

func (s *Server) Shutdown() error {
	_, err := s.engine.Request("Shutdown", nil)
	return err
}

func (s *Server) Uninstall(password string) error {
	// Clean up native host registration (user's HKCU registry keys)
	if err := nativehost.Remove(); err != nil {
		log.Printf("Failed to remove native host: %v", err)
	}
	_, err := s.engine.Request("Uninstall", map[string]string{"password": password})
	return err
}

// --- Settings & History ---

func (s *Server) GetAutostartStatus() (json.RawMessage, error) {
	return s.engine.Request("GetAutostartStatus", nil)
}

func (s *Server) EnableAutostart() error {
	_, err := s.engine.Request("EnableAutostart", nil)
	return err
}

func (s *Server) DisableAutostart() error {
	_, err := s.engine.Request("DisableAutostart", nil)
	return err
}

func (s *Server) ClearAppHistory(password string) error {
	_, err := s.engine.Request("ClearAppHistory", map[string]string{"password": password})
	return err
}

func (s *Server) ClearWebHistory(password string) error {
	_, err := s.engine.Request("ClearWebHistory", map[string]string{"password": password})
	return err
}
