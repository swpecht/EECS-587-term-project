# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  # All Vagrant configuration is done here. The most common configuration
  # options are documented and commented below. For a complete reference,
  # please see the online documentation at vagrantup.com.

  # Every Vagrant virtual environment requires a box to build off of.
  config.vm.box = "base"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # If true, then any SSH connections made will enable agent forwarding.
  # Default value: false
  # config.ssh.forward_agent = true

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  # config.vm.provider "virtualbox" do |vb|
  #   # Don't boot with headless mode
  #   vb.gui = true
  #
  #   # Use VBoxManage to customize the VM. For example to change memory:
  #   vb.customize ["modifyvm", :id, "--memory", "1024"]
  # end
  #
  # View the documentation for the provider you're using for more
  # information on available options.

$script = <<SCRIPT
# Require update call
apt-get update -y

# Prepare python and instalal pip
apt-get install python-dev -y
wget https://bootstrap.pypa.io/get-pip.py
python get-pip.py

# Install fabric and dependencies
pip install ecdsa
pip install paramiko
pip install pycrypto
pip install fabric

# Dev tools
pip install numpy
apt-get install g++ git unzip vim -y

# OpenMPI
apt-get install libopenmpi-dev libopenmpi-dbg -y
pip install mpi4py

# Caffe
apt-get install libprotobuf-dev libleveldb-dev libsnappy-dev libopencv-dev libboost-all-dev libhdf5-serial-dev -y
# glog
wget https://google-glog.googlecode.com/files/glog-0.3.3.tar.gz
tar zxvf glog-0.3.3.tar.gz
cd glog-0.3.3
./configure
make && make install
# gflags
wget https://github.com/schuhschuh/gflags/archive/master.zip
unzip master.zip
cd gflags-master
mkdir build && cd build
export CXXFLAGS="-fPIC" && cmake .. && make VERBOSE=1
make && make install
# lmdb
git clone git://gitorious.org/mdb/mdb.git
cd mdb/libraries/liblmdb
make && make install
SCRIPT

    config.vm.provision "shell", inline: $script

end