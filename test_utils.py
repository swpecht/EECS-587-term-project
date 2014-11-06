import unittest
from utils import *


class TestDataLoading(unittest.TestCase):

    def test_extract_data(self):
        data, labels = extract_data('data/cifar-10-batches-py/data_batch_1')
        # 4d array for the input data
        self.assertEqual(data.ndim, 4)

        # 1d array for the labels
        self.assertEqual(labels.ndim, 1)


if __name__ == '__main__':
    unittest.main()
