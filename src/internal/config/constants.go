package config

const (
	// AppName is used for directories (e.g. in AppData)
	AppName = "VedaAnchor"

	// ServiceName is the Windows Service name
	ServiceName = "VedaAnchorAgent"

	// PipeName is the named pipe address for IPC (Engine)
	PipeName = `\\.\pipe\veda-anchor`

	// AgentPipeName is the named pipe address for Agent UI server
	AgentPipeName = `\\.\pipe\veda-anchor-agent`
)
