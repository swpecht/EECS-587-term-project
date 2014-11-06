from utils import extract_data
import caffe
solver = caffe.SGDSolver('cifar10_quick_solver.prototxt')

data, labels = extract_data('data/cifar-10-batches-py/data_batch_1')

solver.net.set_input_arrays(data, labels)
