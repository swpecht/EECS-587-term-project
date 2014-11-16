package main

import (
	"./communication"
	"fmt"
)

//USE vagrant folders to setup proper package directory / settings for go
//https://golang.org/doc/code.html

func main() {
	client := communication.Client{}
	client.Start(7946)

	client2 := communication.Client{}
	client2.Start(7947)

	var num int

	client2.Join([]string{"0.0.0.0:7946"})
	num = client.NumMembers()
	fmt.Println(num)
}
