package DUP

import (
	"github.com/hashicorp/memberlist"
	"net"
	"strconv"
)

const (
	memberlist_starting_port int = 7946
	tcp_offset               int = 100
)

type ClientFactory struct {
	num_created int
}

func (f *ClientFactory) NewClient() (client, error) {
	c := client{}
	c.ActiveMembers = make(map[string]Node)
	c.pendingMembers = new([]Node)
	//c.EventChannel = &make(chan event)
	// Start event processing
	// go c.processEvents()

	var config *memberlist.Config = memberlist.DefaultLocalConfig()
	config.BindPort = memberlist_starting_port + f.num_created
	c.Name = config.Name + ":" + strconv.Itoa(memberlist_starting_port) + "-" + strconv.Itoa(f.num_created)
	config.Name = c.Name
	config.AdvertisePort = memberlist_starting_port + f.num_created
	config.Events = c
	list, err := memberlist.Create(config)

	c.memberList = list

	c.node = Node{
		Addr: net.ParseIP(config.BindAddr),
		Port: config.BindPort + tcp_offset,
	}

	f.num_created += 1

	return c, err
}
