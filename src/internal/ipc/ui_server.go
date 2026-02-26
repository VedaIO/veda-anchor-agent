//go:build windows

package ipc

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"veda-anchor-agent/src/internal/config"

	"github.com/Microsoft/go-winio"
)

type UIServer struct {
	engineClient *EngineClient
}

func NewUIServer(engineClient *EngineClient) *UIServer {
	return &UIServer{engineClient: engineClient}
}

func (s *UIServer) Start() error {
	address := config.AgentPipeName

	config := &winio.PipeConfig{
		SecurityDescriptor: "D:P(A;;GA;;;AU)",
	}
	listener, err := winio.ListenPipe(address, config)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}
	defer listener.Close()

	log.Printf("[Agent] UI IPC Server listening on %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[Agent] Error accepting connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *UIServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			return
		}

		// Forward to Engine and get response
		result, err := s.forwardToEngine(req.Method, req.Params)
		
		resp := Response{ID: req.ID}
		if err != nil {
			resp.Error = err.Error()
		} else {
			resp.Result = result
		}

		if err := encoder.Encode(resp); err != nil {
			log.Printf("[Agent] Error encoding response: %v", err)
			return
		}
	}
}

func (s *UIServer) forwardToEngine(method string, params json.RawMessage) (interface{}, error) {
	// Re-marshal params as generic interface{}
	var paramsObj interface{}
	if len(params) > 0 {
		json.Unmarshal(params, &paramsObj)
	}

	result, err := s.engineClient.Request(method, paramsObj)
	if err != nil {
		return nil, err
	}
	return result, nil
}
