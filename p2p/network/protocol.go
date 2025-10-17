package network

import (
	"errors"
)

type Protocol interface {
	ProtocolInfo() *ProtocolInfo
	IsVersionSupported(string) (bool, error)
	IsPortSupported(int) (bool, error)
	IsProtocolSupported(string) (bool, error)
	GetDefaultVersion() string
	GetDefaultPort() int
	GetDefaultProtocol() string
}

type ProtocolVersion string

const (
	protocolV1 ProtocolVersion = "1.0.0"
	protocolV2 ProtocolVersion = "1.1.0"
)

type ProtocolPort int

const (
	protocolPortV1 ProtocolPort = 10123
	protocolPortV2 ProtocolPort = 10224
)

type TransportProtocol string

const (
	TCPProtocol TransportProtocol = "tcp"
	UDPProtocol TransportProtocol = "udp"
)

var (
	ErrProtocolVersionNotSupported   = errors.New("protocol version not supported")
	ErrProtocolPortNotSupported      = errors.New("protocol port not supported")
	ErrTransportProtocolNotSupported = errors.New("transport protocol not supported")
)

// ProtcolInfo defines the protocol used in the p2p network
type ProtocolInfo struct {
	TransportProtocol TransportProtocol
	Version           ProtocolVersion
	Port              ProtocolPort
}

// NewProtocolInfo creates a new ProtocolInfo instance
func NewProtocolInfo(transport TransportProtocol) Protocol {
	return &ProtocolInfo{
		TransportProtocol: transport,
		Version:           protocolV2,
		Port:              protocolPortV2,
	}
}

// GetProtocolInfo returns the protocol information
func (pi *ProtocolInfo) ProtocolInfo() *ProtocolInfo {
	return pi
}

// supportedVersions returns a list of supported protocol versions
func supportedVersions() []string {
	return []string{string(protocolV1), string(protocolV2)}
}

// IsVersionSupported checks if the given version is supported
func (pi *ProtocolInfo) IsVersionSupported(version string) (bool, error) {
	for _, v := range supportedVersions() {
		if v == version {
			return true, nil
		}
	}
	return false, ErrProtocolVersionNotSupported
}

// GetDefaultVersion returns the default protocol version
func (pi *ProtocolInfo) GetDefaultVersion() string {
	return string(protocolV2)
}

// supportedPorts returns a list of supported protocol ports
func supportedPorts() []int {
	return []int{int(protocolPortV1), int(protocolPortV2)}
}

// IsPortSupported checks if the given port is supported
func (pi *ProtocolInfo) IsPortSupported(port int) (bool, error) {
	for _, p := range supportedPorts() {
		if p == port {
			return true, nil
		}
	}
	return false, ErrProtocolPortNotSupported
}

// GetDefaultPort returns the default protocol port
func (pi *ProtocolInfo) GetDefaultPort() int {
	return int(protocolPortV2)
}

// supportedProtocols returns a list of supported transport protocols
func supportedProtocols() []string {
	return []string{string(TCPProtocol), string(UDPProtocol)}
}

// IsProtocolSupported checks if the given protocol is supported
func (pi *ProtocolInfo) IsProtocolSupported(proto string) (bool, error) {
	for _, p := range supportedProtocols() {
		if p == proto {
			return true, nil
		}
	}
	return false, ErrTransportProtocolNotSupported
}

// GetDefaultProtocol returns the default transport protocol
func (pi *ProtocolInfo) GetDefaultProtocol() string {
	return string(TCPProtocol)
}
