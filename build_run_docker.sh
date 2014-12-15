#!/bin/bash

# Build the go program
go build -o /vagrant/bin/docker_example $GOPATH/src/github.com/swpecht/EECS-587-term-project/docker.go

# Build the docker image
sudo docker build -t swpecht/broadcast-test ./

# Run the docker image
sudo docker run --net=host -P swpecht/broadcast-test /tmp/docker_example

