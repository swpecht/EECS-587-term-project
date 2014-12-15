from fabric.api import *
from random import randint
import time

# env.use_ssh_config = True
env.user = 'steven.pecht'
# Remote hosts
env.hosts = ['146.148.78.230',
             '146.148.89.88',
             '146.148.45.143']
# env.hosts = ['130.211.122.241']


def build():
    local(
        'go build -o /vagrant/bin/docker_example $GOPATH/src/github.com/swpecht/EECS-587-term-project/docker.go')


def transfer_master():
    build()
    print('Transfering files...')
    put(local_path="./bin/docker_example", remote_path='/home/steven.pecht/')
    run('chmod +x /home/steven.pecht/docker_example')


@parallel
def transfer_test():
    build()
    print('Transfering files...')
    put(local_path="./bin/docker_example", remote_path='/home/steven.pecht/')
    run('chmod +x /home/steven.pecht/docker_example')


@parallel
def run_test():
    print('Sleeping...')
    # Add random delay here to help with parallelizing
    sleep_time = randint(0, 30)
    time.sleep(sleep_time)
    print('Running...')
    run('/home/steven.pecht/docker_example')
