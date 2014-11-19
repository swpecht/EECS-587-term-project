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
	// Initialize variables
	c.ActiveMembers = make(map[string]Node)
	c.pendingMembers = new([]Node)
	c.msgChannel = make(chan Message)
	c.closeChannel = make(chan bool)
	c.barrierChannel = make(chan string)
	c.connectionPool = make(map[string]*net.TCPConn)

	// Start event processing
	go c.startMessageHandling()

	// Configure the MemberList
	var config *memberlist.Config = memberlist.DefaultLocalConfig()
	config.BindPort = memberlist_starting_port + f.num_created
	c.Name = config.Name + ":" + strconv.Itoa(memberlist_starting_port) + "-" + strconv.Itoa(f.num_created)
	config.Name = c.Name
	config.AdvertisePort = memberlist_starting_port + f.num_created
	config.Events = c

	// Configure the local Node data
	c.node = Node{
		Name: c.Name,
		Addr: net.ParseIP(config.BindAddr),
		Port: config.BindPort + tcp_offset,
	}

	tcpAddr := c.node.GetTCPAddr()
	c.tcpListener, err = net.ListenTCP("tcp", &tcpAddr)
	if err != nil {
		return
	}
	// Start the TCP listener
	go tcpListen(c.tcpListener, c.msgChannel)

	list, err := memberlist.Create(config)
	if err != nil {
		return
	}
	c.memberList = list

	f.num_created += 1

	return
}
