package upnp

type Protocol string

const (
	TCP Protocol = "TCP"
	UDP Protocol = "UDP"
)

type MappingPort struct {
	Protocol   Protocol
	Describe   string
	LocalPort  uint16
	RemotePort uint16
}

func NewMappingPort(protocol Protocol, describe string, localPort uint16, remotePort uint16) *MappingPort {
	return &MappingPort{Protocol: protocol, Describe: describe, LocalPort: localPort, RemotePort: remotePort}
}
