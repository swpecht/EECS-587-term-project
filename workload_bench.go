package main

import (
	"fmt"
	"github.com/swpecht/GoMM"
	"time"
)

func main() {
	numIterations := 50
	headNode := "10.0.2.15:7946" // "130.211.122.241:7946"

	f := GoMM.ClientFactory{}
	client, err := GoMM.GetTCPClient(f)
	if err != nil {
		fmt.Println("Failed to create client", err.Error())
		return
	}

	client.Start()
	client.Join(headNode)

	client.WaitActive()
	// Check if first worker
	iter := 0
	if client.NumActiveMembers() > 1 {
		// If not, should be getting the iteration
		msg := <-client.BroadcastChannel
		iter = int(msg.FloatData[0])
	}

	id := client.GetId()
	numNodes := client.NumActiveMembers()
	start := time.Now()
	for i := iter; i < numIterations; i++ {
		oldId := id
		oldNum := numNodes

		UpdatePool(client, i)

		id = client.GetId()
		numNodes = client.NumActiveMembers()

		// Shuffle if update or first iteration
		if oldId != id || oldNum != numNodes || i == iter {
			// Need to reshuffle data
			Shuffle(id, numNodes, client)
		}
		DoIteration(id, numNodes, i, client)
	}
	elapsed := time.Since(start)
	fmt.Println("Benchmark took", elapsed, "for", numIterations-iter, "iterations on", numNodes, "nodes")
	fmt.Println("Average seconds per iteration:", elapsed.Seconds()/float64(numIterations-iter))
}

// Returns the current iteration
func UpdatePool(client *GoMM.Client, iter int) int {
	if client.NumMembers() != client.NumActiveMembers() {
		client.UpdateActiveMembers()
		// If client 0, sent the iteration broadcast to everyone
		if client.GetId() == 0 {
			client.Broadcast([]string{}, []float64{float64(iter)})
		}
		// Receive the iteration broadcast
		msg := <-client.BroadcastChannel
		iter = int(msg.FloatData[0])

	}
	fmt.Println("Doing UpdatePool barrier")
	client.Barrier()
	return iter
}

// Shuffles / distributes the data between all nods
func Shuffle(id, numNodes int, client *GoMM.Client) {
	// Do a barrier for worst case coordination
	fmt.Println("Doing Shuffle barrier")
	client.Barrier()
}

// Perform the actual iterations
func DoIteration(id, numNodes, iter int, client *GoMM.Client) {
	duration := time.Duration(500 / numNodes)
	time.Sleep(duration * time.Millisecond)
	// Assume worst case communications
	// client.Barrier()
	fmt.Printf("Iteration %v \n", iter)
}
