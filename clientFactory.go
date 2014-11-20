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

func (f *ClientFactory) NewClient() (c client, err error) {
	c = client{}
	f.initializeData(&c)
	err = f.initializeTCPListener(&c)
	if err != nil {
		return c, err
	}
	f.startMessageHandling(&c)
	f.startActivateHandling(&c)
	err = f.initializeMemberList(&c)
	if err != nil {
		return c, err
	}

	f.num_created += 1

	return
}

func (f *ClientFactory) initializeData(c *client) {
	// Initialize variables
	c.ActiveMembers = make(map[string]Node)
	c.pendingMembers = new([]Node)
	c.msgIncoming = make(chan Message)
	c.closeChannel = make(chan bool)
	c.barrierChannel = make(chan string)
	c.activateChannel = make(chan Message)

	var config *memberlist.Config = memberlist.DefaultLocalConfig()
	c.Name = config.Name + ":" + strconv.Itoa(memberlist_starting_port) + "-" + strconv.Itoa(f.num_created)

	// Configure the local Node data
	c.node = Node{
		Name: c.Name,
		Addr: net.ParseIP(config.BindAddr),
		Port: config.BindPort + tcp_offset,
	}
}

// Start event processing
func (f *ClientFactory) startMessageHandling(c *client) {
	go c.startMessageHandling()
}

func (f *ClientFactory) startActivateHandling(c *client) {
	go c.startActivateHandling()
}

func (f *ClientFactory) initializeTCPListener(c *client) error {
	// tcpAddr := c.node.GetTCPAddr()
	// var err error
	// // c.tcpListener, err = net.ListenTCP("tcp", &tcpAddr)
	// if err != nil {
	// 	return err
	// }
	// // Start the TCP listener
	// go tcpListen(c.tcpListener, c.msgChannel)
	return nil
}

func (f *ClientFactory) initializeMemberList(c *client) error {
	var config *memberlist.Config = memberlist.DefaultLocalConfig()
	config.BindPort = memberlist_starting_port + f.num_created
	config.Name = c.Name
	config.AdvertisePort = memberlist_starting_port + f.num_created
	config.Events = c

	list, err := memberlist.Create(config)
	c.memberList = list
	return err
}
