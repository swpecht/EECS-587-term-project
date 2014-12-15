package main

import (
	"fmt"
	"github.com/swpecht/GoMM"
	"time"
)

func main() {
	// Test configurations
	numNodes := 2
	numIterations := 10

	// Client data
	isHeadNode := true
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
		isHeadNode = false
		client.Join("10.0.2.15:7946")
		client.WaitActive()
	}

	for client.NumActiveMembers() < numNodes {
		if client.GetNumPendingMembers() > 0 {
			client.UpdateActiveMembers()
		}
		time.Sleep(50 * time.Millisecond)
	}

	if isHeadNode {
		stringData := []string{"Hello", "World"}
		floatData := []float64{2.0, 48182.2}

		start := time.Now()
		for i := 0; i < numIterations; i++ {
			clients[0].Broadcast(stringData, floatData)
			<-client.BroadcastChannel
		}
		elapsed := time.Since(start)
		fmt.Println("Benchmakr took", elapsed, "for", numIterations, "iterations")
		fmt.Println("Average seconds per iteration:", elapsed.Seconds()/float64(numIterations))
	} else {
		<-client.BroadcastChannel
		start := time.Now()
		for i := 0; i < numIterations-2; i++ {
			<-client.BroadcastChannel
		}
		elapsed := time.Since(start)
		fmt.Println("Benchmakr took", elapsed, "for", numIterations-2, "iterations")
		fmt.Println("Average seconds per iteration:", elapsed.Seconds()/float64(numIterations-2))
	}

	// Let all of the messages send
	time.Sleep(1 * time.Second)

}
