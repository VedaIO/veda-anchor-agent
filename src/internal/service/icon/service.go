package icon

import (
	"log"
	"sync"
	"veda-anchor-agent/src/internal/platform/executable"
	get_app_icon "veda-anchor-agent/src/internal/platform/get_app_icon"
)

// AppDetails contains the commercial name and base64 icon of an application.
type AppDetails struct {
	CommercialName string
	IconBase64     string
}

// Service handles app metadata extraction and caching.
type Service struct {
	iconCache   map[string]string
	iconCacheMu sync.Mutex
}

// NewService creates a new icon service.
func NewService() *Service {
	return &Service{
		iconCache: make(map[string]string),
	}
}

// GetAppDetails retrieves the commercial name and icon for a given executable path.
func (s *Service) GetAppDetails(exePath string) AppDetails {
	// Get the commercial name
	commercialName, err := executable.GetCommercialName(exePath)
	if err != nil {
		log.Printf("Failed to get commercial name for %s: %v", exePath, err)
		commercialName = ""
	}

	s.iconCacheMu.Lock()
	iconBase64, ok := s.iconCache[exePath]
	s.iconCacheMu.Unlock()

	if ok {
		return AppDetails{
			CommercialName: commercialName,
			IconBase64:     iconBase64,
		}
	}

	// Extract and cache icon
	iconBase64, err = get_app_icon.GetAppIconAsBase64(exePath)
	if err != nil {
		log.Printf("Failed to get icon for %s: %v", exePath, err)
	}

	s.iconCacheMu.Lock()
	s.iconCache[exePath] = iconBase64
	s.iconCacheMu.Unlock()

	return AppDetails{
		CommercialName: commercialName,
		IconBase64:     iconBase64,
	}
}
