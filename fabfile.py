from fabric.api import *
import time
import random

# env.use_ssh_config = True
env.user = 'steven.pecht'
# Remote hosts
# env.hosts = ['146.148.78.230',
#              '146.148.89.88',
# '146.148.45.143'] #,
# '130.211.122.241']
env.roledefs.update({
    'master': ['130.211.122.241'],
    'slaves': ['108.59.83.231',
               '146.148.88.142',
               '146.148.78.230',
               '146.148.89.88',
               '162.222.183.59',
               '146.148.33.21',
               '146.148.64.85',
               '130.211.114.110',
               '23.236.60.118',
               '146.148.49.115',
               '130.211.123.49',
               '146.148.57.151',
               '146.148.92.239',
               '107.178.214.198',
               '162.222.178.138',
               '162.222.182.29',
               '146.148.62.98',
               '146.148.49.92',
               '130.211.113.31',
               '146.148.65.133',
               # '162.222.177.113',
               # '146.148.45.143',
               ]
})


@runs_once
def build():
    local(
        'go build -o /vagrant/bin/distributed_bench $GOPATH/src/github.com/swpecht/EECS-587-term-project/distributed_bench.go')


@parallel
@roles('master', 'slaves')
def transfer_test():
    print('Transfering files...')
    put(local_path="./bin/distributed_bench",
        remote_path='/home/steven.pecht/')
    run('chmod +x /home/steven.pecht/distributed_bench')
    transfer_multi_run()


@parallel
@roles('master', 'slaves')
def transfer_multi_run():
    put(local_path="./run_test.sh",
        remote_path='/home/steven.pecht/')
    run('chmod +x /home/steven.pecht/run_test.sh')


@parallel
@roles('slaves')
def run_test():
    print('Sleeping...')
    # Add random delay here to help with parallelizing
    # sleep_time = random.randint(0, 30)
    # time.sleep(sleep_time)
    print('Running...')
    run('/home/steven.pecht/distributed_bench')


@parallel
@roles('slaves')
def run_multi_test():
    print('Running...')
    sleep_time = random.randint(0, 30)
    time.sleep(sleep_time)
    run('/home/steven.pecht/run_test.sh')
