package DUP

import (
	"net"
	"testing"
)

func TestGetAddress(t *testing.T) {
	node := Node{
		Addr: net.ParseIP("192.168.1.1"),
		Port: 500,
	}

	if node.GetConnectionAddress() != "192.168.1.1:500" {
		t.Error("Incorrect connection address.")
	}
}
