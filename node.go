package DUP

import (
	"net"
)

type Node struct {
	Addr net.IP // The address this node can be access at
	Port int    // The port this node listens for connections on
}

// Retruns the connection address for this node
func (n Node) GetTCPAddr() net.TCPAddr {
	tcpAddr := net.TCPAddr{IP: n.Addr, Port: n.Port}
	return tcpAddr
}
