//go:build windows

package ipc

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"veda-anchor-agent/src/internal/config"

	"github.com/Microsoft/go-winio"
)

type EngineClient struct {
	conn net.Conn
}

func NewEngineClient() (*EngineClient, error) {
	address := config.PipeName
	timeout := 5 * time.Second

	conn, err := winio.DialPipe(address, &timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Engine: %w", err)
	}

	log.Printf("[Agent] Connected to Engine at %s", address)
	return &EngineClient{conn: conn}, nil
}

func (c *EngineClient) Request(method string, params interface{}) (json.RawMessage, error) {
	id := generateID()
	paramsJSON, _ := json.Marshal(params)

	req := Request{
		ID:     id,
		Method: method,
		Params: paramsJSON,
	}

	encoder := json.NewEncoder(c.conn)
	if err := encoder.Encode(req); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	decoder := json.NewDecoder(c.conn)
	var resp Response
	if err := decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.ID != id {
		return nil, fmt.Errorf("request ID mismatch")
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("Engine error: %s", resp.Error)
	}

	if resp.Result == nil {
		return nil, nil
	}

	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}
	return resultBytes, nil
}

func (c *EngineClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
