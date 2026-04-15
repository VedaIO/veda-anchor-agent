package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"veda-anchor-agent/src/internal/ipc"
	"veda-anchor-agent/src/internal/service/icon"
)

type Server struct {
	engine          *ipc.EngineClient
	icons           *icon.Service
	IsAuthenticated bool
	Mu              sync.Mutex
}

func NewServer(engine *ipc.EngineClient) *Server {
	return &Server{
		engine: engine,
		icons:  icon.NewService(),
	}
}

type AppDetailsResponse struct {
	CommercialName string `json:"commercialName"`
	Icon           string `json:"icon"`
}

func (s *Server) GetAppDetails(exePath string) (AppDetailsResponse, error) {
	details := s.icons.GetAppDetails(exePath)
	return AppDetailsResponse{
		CommercialName: details.CommercialName,
		Icon:           details.IconBase64,
	}, nil
}

func (s *Server) GetWebDetails(domain string) (json.RawMessage, error) {
	return s.engine.Request("GetWebDetails", map[string]string{"domain": domain})
}

func (s *Server) CheckChromeExtension() bool {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return false
	}
	heartbeatPath := filepath.Join(cacheDir, "VedaAnchor", "extension_heartbeat")
	content, err := os.ReadFile(heartbeatPath)
	if err != nil {
		return false
	}
	var lastPing int64
	if _, err := fmt.Sscanf(string(content), "%d", &lastPing); err != nil {
		return false
	}
	return time.Since(time.Unix(lastPing, 0)) < 10*time.Second
}
