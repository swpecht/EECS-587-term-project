package main

import (
	"fmt"
	"github.com/swpecht/GoMM"
	"time"
)

func main() {
	// Test configurations
	numNodes := 7
	numIterations := 10
	headNode := "130.211.122.241:7946"

	// Client data
	isHeadNode := true
	clients, err := GoMM.GetTCPClients(1)

	if err != nil {
		fmt.Printf("Failed to create client: %s", err.Error())
		return
	}
	client := clients[0]
	fmt.Println("Created client on:", client.JoinAddr())

	client.Start()
	if client.JoinAddr() == "10.240.94.200:7946" {
		fmt.Println("This node is head node")
	} else {
		fmt.Println("This node is NOT head node")
		isHeadNode = false
		client.Join(headNode)
		client.WaitActive()
	}

	for client.NumActiveMembers() < numNodes {
		if client.NumMembers() == numNodes {
			client.UpdateActiveMembers()
		}
		time.Sleep(50 * time.Millisecond)
	}

	start := time.Now()
	if isHeadNode {
		// Broacast the message with the head node id
		stringData := []string{"Hello", "World"}
		floatData := []float64{float64(client.GetId())}

		// Broadcast the messages
		for i := 0; i < numIterations; i++ {
			clients[0].Broadcast(stringData, floatData)
			<-client.BroadcastChannel
		}

	} else {
		// Receive messages
		for i := 0; i < numIterations; i++ {
			<-client.BroadcastChannel
		}
	}
	// Barrier for timing
	client.Barrier()
	elapsed := time.Since(start)
	fmt.Println("Benchmark took", elapsed, "for", numIterations, "iterations on", numNodes, "nodes")
	fmt.Println("Average seconds per iteration:", elapsed.Seconds()/float64(numIterations-2))

	// Let all of the messages send, the entire network is required to send messages
	time.Sleep(40 * time.Second)

}
