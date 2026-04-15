package executable

import (
	"path/filepath"
	"strings"

	"github.com/bi-zone/go-fileversion"
)

// GetCommercialName retrieves the commercial name of the application.
// It tries to find the FileDescription, ProductName, or OriginalFilename from the executable's version info.
// If none are found, it falls back to the filename without extension.
func GetCommercialName(exePath string) (string, error) {
	info, err := fileversion.New(exePath)
	var commercialName string
	if err == nil {
		commercialName = info.FileDescription()
		if commercialName == "" {
			commercialName = info.ProductName()
		}
		if commercialName == "" {
			commercialName = info.OriginalFilename()
		}
	}

	// Fallback to filename
	if commercialName == "" {
		commercialName = strings.TrimSuffix(filepath.Base(exePath), filepath.Ext(exePath))
	}

	return commercialName, nil
}
