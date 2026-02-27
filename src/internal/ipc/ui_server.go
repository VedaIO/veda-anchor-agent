//go:build windows

package ipc

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"syscall"
	"unsafe"

	"veda-anchor-agent/src/internal/config"

	"github.com/Microsoft/go-winio"
	"golang.org/x/sys/windows"
)

type UIServer struct {
	engineClient *EngineClient
	uiExePath    string
}

func NewUIServer(engineClient *EngineClient) *UIServer {
	uiPath := filepath.Join(config.ProgramFiles(), "VedaAnchor", "veda-anchor-ui.exe")

	return &UIServer{
		engineClient: engineClient,
		uiExePath:    uiPath,
	}
}

func (s *UIServer) Start() error {
	address := config.AgentPipeName

	pc := &winio.PipeConfig{}
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

		if !s.validateClient(conn) {
			log.Printf("[Agent] Rejected connection from unverified client")
			conn.Close()
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *UIServer) validateClient(conn net.Conn) bool {
	connVal, ok := conn.(interface{ Handle() windows.Handle })
	if !ok {
		log.Printf("[Agent] Cannot get handle from connection type")
		return false
	}

	handle := connVal.Handle()
	var pid uint32
	ret, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("GetNamedPipeClientProcessId").Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&pid)),
	)
	if ret == 0 {
		log.Printf("[Agent] GetNamedPipeClientProcessId failed")
		return false
	}

	exePath, err := queryProcessPath(pid)
	if err != nil {
		log.Printf("[Agent] Failed to query client process path: %v", err)
		return false
	}

	log.Printf("[Agent] Client connected from: %s", exePath)
	return exePath == s.uiExePath
}

func queryProcessPath(pid uint32) (string, error) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(h)

	var pathBuf [windows.MAX_PATH]uint16
	pathLen := uint32(len(pathBuf))
	err = windows.QueryFullProcessImageName(h, 0, &pathBuf[0], &pathLen)
	if err != nil {
		return "", err
	}

	return windows.UTF16ToString(pathBuf[:pathLen]), nil
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
