package DUP

import (
	"strconv"
)

type Node struct {
	Addr string // The address this node can be access at
	Port int    // The port this node listens for connections on
}

// Retruns the connection address for this node
func (n Node) GetConnectionAddress() string {
	return n.Addr + strconv.Itoa(n.Port)
}
