import numpy as np


def unpickle(file):
    import cPickle
    fo = open(file, 'rb')
    dict = cPickle.load(fo)
    fo.close()
    return dict


def extract_data(file, channels=3, width=32, height=32):
    """ Returns the [4-d numpy data array, labels] of the target file.
    The first dimensions of the data array are:
        1. The number of pictures
        2. Color channel (3)
        3. Pixel value (1024, 32x32 image)
        4. Place holder to meet caffe 4d requirement
    The data array has type float32
        Ref on data shape: https://groups.google.com/forum/#!msg/caffe-users/
                            wWSGX4vmAh4/ivEjy-pPLckJ
        """
    raw_data = unpickle(file)
    labels = np.array(raw_data['labels'])
    raw_data = np.array(raw_data['data'])

    data = np.zeros((10000, channels, height, width), dtype=np.uint8)

    for channel in range(0, channels):
        for row in range(0, height):
            channel_start = 1024 * channel + width * row
            data[:, channel, row, :] = raw_data[
                :, channel_start:channel_start + width]

    # Convert to floats
    data = data.astype(np.float32)
    labels = labels.astype(np.float32)
    return [data, labels]
