package DUP

import (
	"net"
	"strconv"
)

type Node struct {
	Addr net.IP // The address this node can be access at
	Port int    // The port this node listens for connections on
}

// Retruns the connection address for this node
func (n Node) GetConnectionAddress() string {
	return n.Addr.String() + strconv.Itoa(n.Port)
}
