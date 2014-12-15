package main

import (
	"fmt"
	"github.com/swpecht/GoMM"
	"time"
)

func main() {
	numNodes := 3

	clients, err := GoMM.GetTCPClients(1)
	if err != nil {
		fmt.Printf("Failed to create client: %s", err.Error())
		return
	}
	client := clients[0]

	client.Start()
	if client.JoinAddr() == "10.0.2.15:7946" {
		fmt.Printf("This node is head node")
	} else {
		fmt.Printf("This node is NOT head node")
		client.Join("10.0.2.15:7946")
		client.WaitActive()
	}

	for client.NumActiveMembers() < numNodes {
		if client.GetNumPendingMembers() > 0 {
			client.UpdateActiveMembers()
		}
		time.Sleep(50 * time.Millisecond)
	}
}
