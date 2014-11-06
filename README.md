EECS-587-term-project
=====================

Distributed image classification and model training.

# Summary

The aim of this project is to simulate a system that classifies images submitted by users. The cifar-10 data set will be used as sample data. Outside of classifying the images, the system should:

* Periodically update the classification models based on new data while still serving new requests
* Take a user's history into account when making a classification of an image

The goals are to:

* Determine how different communication and model update freuqencies impact system performance over time
* Develop a method to keep user models updated despite the large number of users
* Examine different approaches to parallelize the training of the models

This project is different from previous work such as the system used by [Google's priority inbox](http://static.googleusercontent.com/media/research.google.com/en/us/pubs/archive/36955.pdf) in that the models used here are not able to be trained in parallel and later combined.

# Implementation

## Environment
The code is designed to be run on [Flux](http://arc.research.umich.edu/flux-and-other-hpc-resources/flux/)

Caffe can be added to the python path using:

```
export PYTHONPATH=/home/software/rhel6/caffe/0.9999/distribute/python:$PYTHONPATH
```

Similalry, protobuf can be added to the path using:

```
export PYTHONPATH=/home/swpecht/term-project/proto/protobuf-2.5.0/:$PYTHONPATH
```

The required modules can be found in the [required_modules](required_modules) file.

# Data collection

# Plan

1. Run cifar10 caffe example to ensure everything works - DONE
2. Run cifar10 caffe example in python
3. Compare training time on various training datasets sizes to compute theoretical possible speedup


# References

* Data desription: http://www.cs.toronto.edu/~kriz/cifar.html
* Data: http://www.cs.toronto.edu/~kriz/cifar-10-python.tar.gz
* Caffe tutorial: http://caffe.berkeleyvision.org/gathered/examples/cifar10.html
* PARALLELIZING NEURAL NETWORK TRAINING FOR CLUSTER SYSTEMS: http://www.cs.swarthmore.edu/~newhall/papers/pdcn08.pdf
* Levels of parallelism: http://techblog.netflix.com/2014/02/distributed-neural-networks-with-gpus.html
* http://stackoverflow.com/questions/5679008/a-i-how-would-i-train-a-neural-network-across-multiple-machines

