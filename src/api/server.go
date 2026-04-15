package api

import (
	"encoding/json"
	"sync"
	"veda-anchor-agent/src/internal/ipc"
)

type Server struct {
	engine          *ipc.EngineClient
	IsAuthenticated bool
	Mu              sync.Mutex
}

func NewServer(engine *ipc.EngineClient) *Server {
	return &Server{engine: engine}
}

type AppDetailsResponse struct {
	CommercialName string `json:"commercialName"`
	Icon           string `json:"icon"`
}

func (s *Server) GetAppDetails(exePath string) (interface{}, error) {
	return s.engine.Request("GetAppDetails", map[string]string{"exePath": exePath})
}

func (s *Server) GetWebDetails(domain string) (json.RawMessage, error) {
	return s.engine.Request("GetWebDetails", map[string]string{"domain": domain})
}
