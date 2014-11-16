EECS-587-term-project
=====================

Distributed image classification and model training.

# Summary
DUP is a decetranlaized messaging platform that supports dynamic communication pool resizing. This library aims to provide a subset of MPIs functionality. The frameworks provides an a solution for problems that do not have geometry related communication.

The membership functionality will be implemented using [Memberlist](https://github.com/hashicorp/memberlist)

Functionality:

* All-to-all communication: Scatter, reduce
* Controlled adding and removing of nodes from the communication pool.

# Implementation

There will be 2 types of membership for the library:

* Active nodes - nodes that are currently involved in communication
* Online nodes - nodes that are online and and member of the Memberlist, but are not involved in communication

The tiered membership system is requried to ensure that collective communication calls are not impacted by nodes joining when other nodes are not expecting it.

## Environment
The library will be implemented in Go.

# Data collection

# Plan

1. Create intial binding code for Memberlist - DONE
2. Create factory for DUP clients
3. Implement broadcast

# References

* Data desription: http://www.cs.toronto.edu/~kriz/cifar.html
* Data: http://www.cs.toronto.edu/~kriz/cifar-10-python.tar.gz
* Caffe tutorial: http://caffe.berkeleyvision.org/gathered/examples/cifar10.html
* PARALLELIZING NEURAL NETWORK TRAINING FOR CLUSTER SYSTEMS: http://www.cs.swarthmore.edu/~newhall/papers/pdcn08.pdf
* Levels of parallelism: http://techblog.netflix.com/2014/02/distributed-neural-networks-with-gpus.html
* http://stackoverflow.com/questions/5679008/a-i-how-would-i-train-a-neural-network-across-multiple-machines

