package DUP

import (
	"github.com/hashicorp/memberlist"
	"strconv"
)

const (
	memberlist_starting_port int = 7946
	tcp_starting_port        int = 8000
)

type ClientFactory struct {
	num_created int
}

func (f *ClientFactory) NewClient() (client, error) {
	c := client{}

	var config *memberlist.Config = memberlist.DefaultLocalConfig()
	config.BindPort = memberlist_starting_port + f.num_created
	c.Name = config.Name + ":" + strconv.Itoa(memberlist_starting_port) + "-" + strconv.Itoa(f.num_created)
	config.Name = c.Name
	config.AdvertisePort = memberlist_starting_port + f.num_created
	config.Events = c
	list, err := memberlist.Create(config)

	c.membersList = list

	c.node = Node{
		Addr: config.BindAddr,
		Port: tcp_starting_port + f.num_created,
	}

	f.num_created += 1

	return c, err
}
