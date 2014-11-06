from fabric.api import *

env.user = 'swpecht'
env.hosts = 'flux-xfer.engin.umich.edu'


def transfer():
    print('Transfering files...')
    # Need to use the flux transfer host
    put(local_path='./*.py', remote_path='~/term-project/')


def get_output(folder):
    print('Downloading files...')
    get(local_path='output/', remote_path='~/' + folder + '/*.o*')

def get_modules():
    print('Downloading module file...')
    get(local_path='required_modules', remote_path='~/privatemodules/default')
