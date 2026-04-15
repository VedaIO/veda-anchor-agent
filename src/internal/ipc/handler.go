package ipc

import "encoding/json"

// RequestHandler processes incoming IPC requests.
// Implemented by api.Server to break the circular dependency.
type RequestHandler interface {
	Dispatch(method string, params json.RawMessage) Response
}
