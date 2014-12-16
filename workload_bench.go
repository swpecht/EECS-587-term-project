package main

import (
	"fmt"
	"github.com/swpecht/GoMM"
)

func main() {
	numIterations := 500
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

	for i := 0; i < numIterations; i++ {
		id, numNodes := UpdatePool(client)
		Shuffle(id, numNodes, client)
		DoIteration(id)
	}
}

// Returns the id, numActiveNodes
func UpdatePool(client *GoMM.Client) (int, int) {
	return 0, 1
}

// Shuffles / distributes the data between all nods
func Shuffle(id, numNodes int, client *GoMM.Client) {

}

// Perform the actual iterations
func DoIteration(id int) {

}
