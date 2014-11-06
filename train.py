from utils import extract_data
import numpy as np
import caffe

solver = caffe.SGDSolver('cifar10_quick_solver.prototxt')

data, labels = extract_data('data/cifar-10-batches-py/data_batch_1')
data_2, labels_2 = extract_data('data/cifar-10-batches-py/data_batch_2')

data = np.concatenate((data, data_2), axis=0)
labels = np.concatenate((labels, labels_2), axis=0)

print labels.shape[0]

solver.net.set_input_arrays(data, labels)

# test_data, test_labels = extract_data('data/cifar-10-batches-py/test_batch')
# solver.test_nets.set_input_arrays(test_data, test_labels)

solver.solve()

print solver.net.params['weights']
