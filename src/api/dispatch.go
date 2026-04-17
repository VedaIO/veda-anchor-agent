package api

import (
	"encoding/json"
	"veda-anchor-agent/src/internal/ipc"
	"veda-anchor-agent/src/internal/tracking"
)

func (s *Server) Dispatch(method string, params json.RawMessage) ipc.Response {
	var result any
	var err error

	switch method {

	// --- Locally handled methods ---

	case "GetAppDetails":
		var p struct {
			ExePath string `json:"exePath"`
		}
		json.Unmarshal(params, &p)
		result, err = s.GetAppDetails(p.ExePath)

	case "CheckChromeExtension":
		result = s.CheckChromeExtension()

	case "GetIsAuthenticated":
		result = s.GetIsAuthenticated()

	case "Logout":
		s.Logout()
		result = true

	// --- Methods that forward to engine but also update local state ---

	case "Login":
		var p struct {
			Password string `json:"password"`
		}
		json.Unmarshal(params, &p)
		result, err = s.Login(p.Password)

	case "SetPassword":
		var p struct {
			Password string `json:"password"`
		}
		json.Unmarshal(params, &p)
		err = s.SetPassword(p.Password)

	case "ClearAppHistory":
		var p struct {
			Password string `json:"password"`
		}
		json.Unmarshal(params, &p)
		err = s.ClearAppHistory(p.Password)
		if err == nil {
			tracking.Reset()
		}

	case "GetAppLeaderboard":
		var p struct {
			Since string `json:"since"`
			Until string `json:"until"`
		}
		json.Unmarshal(params, &p)
		raw, fwdErr := s.GetAppLeaderboard(p.Since, p.Until)
		if fwdErr != nil {
			return ipc.Response{Error: fwdErr.Error()}
		}
		// Enrich icons
		var items []AppLeaderboardItem
		if err := json.Unmarshal(raw, &items); err == nil {
			for i := range items {
				details := s.icons.GetAppDetails(items[i].ExecutablePath)
				if details.CommercialName != "" {
					items[i].Name = details.CommercialName
				}
				items[i].Icon = details.IconBase64
			}
			enriched, _ := json.Marshal(items)
			return ipc.Response{Result: json.RawMessage(enriched)}
		}
		return ipc.Response{Result: raw}

	case "GetScreenTime":
		raw, fwdErr := s.GetScreenTime()
		if fwdErr != nil {
			return ipc.Response{Error: fwdErr.Error()}
		}
		// Enrich icons
		var items []ScreenTimeItem
		if err := json.Unmarshal(raw, &items); err == nil {
			for i := range items {
				details := s.icons.GetAppDetails(items[i].ExecutablePath)
				if details.CommercialName != "" {
					items[i].Name = details.CommercialName
				}
				items[i].Icon = details.IconBase64
			}
			enriched, _ := json.Marshal(items)
			return ipc.Response{Result: json.RawMessage(enriched)}
		}
		return ipc.Response{Result: raw}

	// --- Everything else: forward to engine ---

	default:
		raw, fwdErr := s.forwardToEngine(method, params)
		if fwdErr != nil {
			return ipc.Response{Error: fwdErr.Error()}
		}
		// Return raw JSON directly to avoid double-encoding
		return ipc.Response{Result: raw}
	}

	if err != nil {
		return ipc.Response{Error: err.Error()}
	}
	return ipc.Response{Result: result}
}

// forwardToEngine sends a raw request to the engine and returns the result
func (s *Server) forwardToEngine(method string, params json.RawMessage) (json.RawMessage, error) {
	var paramsObj any
	if len(params) > 0 {
		json.Unmarshal(params, &paramsObj)
	}
	return s.engine.Request(method, paramsObj)
}
