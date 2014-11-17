package main

import (
	"./communication"
	"fmt"
)

//USE vagrant folders to setup proper package directory / settings for go
//https://golang.org/doc/code.html

func main() {
	factory := communication.ClientFactory{}
	client, _ := factory.NewClient()

	client2, _ := factory.NewClient()
	client3, _ := factory.NewClient()

	client2.Join([]string{"0.0.0.0:7946"})
	client3.Join([]string{"0.0.0.0:7946"})
	num := client.NumMembers()
	fmt.Println(num)
}
