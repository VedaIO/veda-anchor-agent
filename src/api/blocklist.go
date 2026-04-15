package api

import (
	"encoding/json"
)

// --- App Blocklist ---

func (s *Server) BlockApps(names []string) error {
	_, err := s.engine.Request("BlockApps", names)
	return err
}

func (s *Server) UnblockApps(names []string) error {
	_, err := s.engine.Request("UnblockApps", names)
	return err
}

func (s *Server) GetAppBlocklist() (json.RawMessage, error) {
	return s.engine.Request("GetAppBlocklist", nil)
}

func (s *Server) ClearAppBlocklist() error {
	_, err := s.engine.Request("ClearAppBlocklist", nil)
	return err
}

func (s *Server) SaveAppBlocklist() (json.RawMessage, error) {
	return s.engine.Request("SaveAppBlocklist", nil)
}

func (s *Server) LoadAppBlocklist(content []byte) error {
	_, err := s.engine.Request("LoadAppBlocklist", content)
	return err
}

// --- Web Blocklist ---

func (s *Server) GetWebBlocklist() (json.RawMessage, error) {
	return s.engine.Request("GetWebBlocklist", nil)
}

func (s *Server) AddWebBlocklist(domain string) error {
	_, err := s.engine.Request("AddWebBlocklist", domain)
	return err
}

func (s *Server) RemoveWebBlocklist(domain string) error {
	_, err := s.engine.Request("RemoveWebBlocklist", domain)
	return err
}

func (s *Server) ClearWebBlocklist() error {
	_, err := s.engine.Request("ClearWebBlocklist", nil)
	return err
}

func (s *Server) SaveWebBlocklist() (json.RawMessage, error) {
	return s.engine.Request("SaveWebBlocklist", nil)
}

func (s *Server) LoadWebBlocklist(content []byte) error {
	_, err := s.engine.Request("LoadWebBlocklist", content)
	return err
}
