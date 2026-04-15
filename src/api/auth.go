package api

import (
	"encoding/json"
)

// GetIsAuthenticated checks if the user is authenticated.
func (s *Server) GetIsAuthenticated() bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	return s.IsAuthenticated
}

// Logout handles the user logout.
func (s *Server) Logout() {
	s.Mu.Lock()
	s.IsAuthenticated = false
	s.Mu.Unlock()
}

// HasPassword checks if a password has been set for the application.
func (s *Server) HasPassword() (json.RawMessage, error) {
	return s.engine.Request("HasPassword", nil)
}

// Login handles the user login.
func (s *Server) Login(password string) (json.RawMessage, error) {
	result, err := s.engine.Request("Login", map[string]string{"password": password})
	if err != nil {
		return nil, err
	}

	// If login succeeded, set local auth state
	var success bool
	json.Unmarshal(result, &success)
	if success {
		s.Mu.Lock()
		s.IsAuthenticated = true
		s.Mu.Unlock()
	}
	return result, nil
}

// SetPassword handles the initial password setup.
func (s *Server) SetPassword(password string) error {
	_, err := s.engine.Request("SetPassword", map[string]string{"password": password})
	if err == nil {
		s.Mu.Lock()
		s.IsAuthenticated = true
		s.Mu.Unlock()
	}
	return err
}
