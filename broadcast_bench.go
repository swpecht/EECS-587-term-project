package main

import (
	"fmt"
	"github.com/swpecht/GoMM"
	"time"
)

func main() {
	numNodes := 10

	factory := GoMM.ClientFactory{}
	headName := "0.0.0.0:7946"
	clients := GoMM.GetLocalClients(numNodes, headName)

	for i := range clients {
		clients[i].Start()
	}

	// Don't have the root node join itsself
	for i := 1; i < len(clients); i++ {
		clients[i].Join(headName)
	}

	// Activate the pending memebers
	clients[0].UpdateActiveMembers()

	// Wait for all members to activate
	for i := 0; i < len(clients); i++ {
		clients[i].WaitActive()
	}

	start := time.Now()
	for i := 0; i < numNodes; i++ {
		fmt.Println("Hello, world.\n", factory, "\n")
		// Run a barrier, maybe have a client output to a channel when it
		// recieves the message, will allow for timing.
		// May require some knowledge of which one will get the message last.
	}
	elapsed := time.Since(start)
	fmt.Println("Benchmakr took", elapsed)

}
