package main

import (
	"fmt"
	"github.com/swpecht/GoMM"
	"time"
)

func main() {
	numNodes := 50
	numIterations := 100

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

	stringData := []string{"Hello", "World"}
	floatData := []float64{2.0, 48182.2}

	start := time.Now()
	for i := 0; i < numIterations; i++ {
		clients[0].Broadcast(stringData, floatData)
		ReceiveAllMessages(clients)
	}
	elapsed := time.Since(start)
	fmt.Println("Benchmakr took", elapsed, "for", numIterations, "iterations")
	fmt.Println("Average seconds per iteration:", elapsed.Seconds()/float64(numIterations))
}

// Receive messages all on clients but the root node
func ReceiveAllMessages(clients []GoMM.Client) {
	for i := 0; i < len(clients); i++ {
		<-clients[i].BroadcastChannel
	}
}
