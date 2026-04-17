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
	handler RequestHandler
}

func NewUIServer(handler RequestHandler) *UIServer {
	return &UIServer{
		handler: handler,
	}
}

func (s *UIServer) Start() error {
	address := config.AgentPipeName

	// Security: Use DACL to restrict pipe access at the kernel level.
	// SY = SYSTEM, BA = Built-in Administrators, AU = Authenticated Users.
	// This is the industry-standard approach — the kernel enforces who can connect,
	// rather than fragile process-path validation after the fact.
	pc := &winio.PipeConfig{
		SecurityDescriptor: "D:(A;;GA;;;SY)(A;;GA;;;BA)(A;;GA;;;AU)",
	}
	listener, err := winio.ListenPipe(address, pc)
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

	log.Printf("[Agent] Accepted UI connection")

	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			return
		}

		resp := s.handler.Dispatch(req.Method, req.Params)
		resp.ID = req.ID

		if err := encoder.Encode(resp); err != nil {
			log.Printf("[Agent] Error encoding response: %v", err)
			return
		}
	}
}
