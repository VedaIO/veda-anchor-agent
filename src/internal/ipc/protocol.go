package ipc

import (
	"encoding/json"
	"veda-anchor-agent/src/internal/config"
)

// Request is a message received from the client.
type Request struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

// Response is a message sent to the client.
type Response struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// GetIPCAddress returns the Engine's pipe address.
func GetIPCAddress() string {
	return config.PipeName
}

// GetAgentIPCAddress returns the Agent's pipe address for UI connections.
func GetAgentIPCAddress() string {
	return config.AgentPipeName
}
