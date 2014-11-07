EECS-587-term-project
=====================

Distributed image classification and model training.

# Summary
There are a number of different levels of parallelism possible when using neural networks. A more detailed analysis can be found on [this netflix blog](http://techblog.netflix.com/2014/02/distributed-neural-networks-with-gpus.html). The main three levels are:

1. Across region
2. Hyperparameter optimization
3. Model training

The first two levels are frequently done in parallel. Parallelizing across regions requires little to no communication. And hyperparamemter optimization can be performed using a number of parallel optimization algorithms.

Training models in parallel poses three main challenges:

* Stopping criteria
* Distributing training data
* Balancing batch size
* Requirements for synchronization

This project builds off of the work of [Dahl et al](http://www.cs.swarthmore.edu/~newhall/papers/pdcn08.pdf). And addresses the following

As well as accounting for:
* Dynamically adjusting the number of nodes
* Testing the methodology on larger sets of ANN problems

# Implementation

Batch size determines training time

## Environment
The code is designed to be run on [Flux](http://arc.research.umich.edu/flux-and-other-hpc-resources/flux/)

Caffe can be added to the python path using:

```bash
export PYTHONPATH=/home/software/rhel6/caffe/0.9999/distribute/python:$PYTHONPATH
```

Similalry, protobuf can be added to the path using:

```bash
export PYTHONPATH=/home/swpecht/term-project/proto/protobuf-2.5.0/:$PYTHONPATH
```

The required modules can be found in the [required_modules](required_modules) file.

# Data collection

# Plan

1. Run cifar10 caffe example to ensure everything works - DONE
2. Run cifar10 caffe example in python - DONE
3. Compare training time on various training datasets sizes to compute theoretical possible speedup


# References

* Data desription: http://www.cs.toronto.edu/~kriz/cifar.html
* Data: http://www.cs.toronto.edu/~kriz/cifar-10-python.tar.gz
* Caffe tutorial: http://caffe.berkeleyvision.org/gathered/examples/cifar10.html
* PARALLELIZING NEURAL NETWORK TRAINING FOR CLUSTER SYSTEMS: http://www.cs.swarthmore.edu/~newhall/papers/pdcn08.pdf
* Levels of parallelism: http://techblog.netflix.com/2014/02/distributed-neural-networks-with-gpus.html
* http://stackoverflow.com/questions/5679008/a-i-how-would-i-train-a-neural-network-across-multiple-machines

