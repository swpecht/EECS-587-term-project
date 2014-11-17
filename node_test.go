package DUP

import (
	"testing"
)

func TestGetAddress(t *testing.T) {
	node := Node{
		Addr: "localhost",
		Port: 500,
	}

	if node.GetConnectionAddress() != "localhost:500" {
		t.Error("Incorrect connection address.")
	}
}
