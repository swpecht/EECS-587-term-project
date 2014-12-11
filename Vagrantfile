# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

$script_go = <<SCRIPT
# Require update call
apt-get update


# Install dev tools
apt-get install -y
                    wget \
                    vim \
                    git

# Install Go
wget https://storage.googleapis.com/golang/go1.3.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.3.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
echo 'export GOPATH=$HOME/go' >> /etc/profile
echo 'export PATH=$PATH:$GOPATH/bin' >> /etc/profile
echo 'export GOBIN=$GOPATH/bin'
mkdir $GOPATH
mkdir $GOPATH/pkg
mkdir $GOPATH/bin
mkdir $GOPATH/src

# Setup go dev environtment
go get github.com/stretchr/testify

SCRIPT

  config.vm.define "go_dev" do |go_dev|
    go_dev.vm.box = "ubuntu/trusty64"
    go_dev.vm.provision "shell", inline: $script_go
    go_dev.vm.network "forwarded_port", guest: 8080, host: 8080

    # Make sure everything is in the right place for go
    config.vm.synced_folder "./", "/home/vagrant/go/src/github.com/swpecht/EECS-587-term-project"
    config.vm.synced_folder "./libraries/GoMM", "/home/vagrant/go/src/github.com/swpecht/GoMM"
    go_dev.vm.provider "virtualbox" do |v|
      v.memory = 2048
    end
  end


end
