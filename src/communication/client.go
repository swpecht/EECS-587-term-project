package communication

import (
	"github.com/hashicorp/memberlist"
	"strconv"
)

// https://github.com/hashicorp/memberlist

type Client struct {
	membersList *memberlist.Memberlist
}

func (c Client) NumMembers() int {
	return c.membersList.NumMembers()
}

func (c *Client) Start(port int) {
	var config *memberlist.Config = memberlist.DefaultLocalConfig()
	config.BindPort = port
	config.Name = strconv.Itoa(port)
	config.AdvertisePort = port
	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	c.membersList = list
}

func (c Client) Join(addresses []string) int {
	n, err := c.membersList.Join(addresses)
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}

	return n
}

func (c Client) GetMembers() string {
	return "MEMBERS2"
}
