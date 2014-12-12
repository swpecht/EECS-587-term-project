package main

import (
	"fmt"
	"github.com/swpecht/GoMM"
	"time"
)

func main() {
	numNodes := 2
	numIterations := 1

	clients := GoMM.GetLocalClients(numNodes)
	headName := clients[0].JoinAddr()

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

	// Get new headname
	headName = clients[1].JoinAddr()
	start := time.Now()
	for i := 0; i < numIterations; i++ {
		clients[0].Close()
		WaitNodeLeave(clients, numNodes)
	}
	elapsed := time.Since(start)
	fmt.Println("Benchmakr took", elapsed, "for", numIterations, "iterations")
	fmt.Println("Average seconds per iteration:", elapsed.Seconds()/float64(numIterations))
}

// Receive messages all on clients but the root node
func WaitNodeLeave(clients []GoMM.Client, activeNodes int) {
	// Don't test 1 since it left
	for i := 1; i < len(clients); i++ {
		for clients[i].NumActiveMembers() == activeNodes {
			// keep looping here
			time.Sleep(time.Millisecond * 10)
		}
	}
}
