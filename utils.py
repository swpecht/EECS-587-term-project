import numpy as np


def unpickle(file):
    import cPickle
    fo = open(file, 'rb')
    dict = cPickle.load(fo)
    fo.close()
    return dict


def extract_data(file):
    """ Returns the [4-d numpy data array, labels] of the target file """
    raw_data = unpickle(file)
    labels = np.array(raw_data['labels'])
    raw_data = np.array(raw_data['data'])

    data = np.zeros((10000, 3, 1024),dtype=np.uint8)
    # Copy data channels
    data[:, 0, :] = raw_data[:, 0:1024]
    data[:, 1, :] = raw_data[:, 1024:2048]
    data[:, 2, :] = raw_data[:, 2048:3072]
    # Make the array 4d for use with caffe
    data = np.expand_dims(data, axis=3)

    print data[0, :]
    return [data, labels]
