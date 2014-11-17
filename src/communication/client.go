package communication

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"strconv"
)

// https://github.com/hashicorp/memberlist

type client struct {
	membersList   *memberlist.Memberlist
	pendingMember []memberlist.Node
	ActiveMembers []memberlist.Node
	Name          string
}

type ClientFactory struct {
	num_created int
}

func (f *ClientFactory) NewClient() (client, error) {
	c := client{}

	var config *memberlist.Config = memberlist.DefaultLocalConfig()
	config.BindPort = 7946 + f.num_created
	c.Name = config.Name + ":" + strconv.Itoa(7946) + "-" + strconv.Itoa(f.num_created)
	config.Name = c.Name
	config.AdvertisePort = 7946 + f.num_created
	config.Events = c
	list, err := memberlist.Create(config)

	c.membersList = list
	f.num_created += 1
	return c, err
}

func (c client) NotifyJoin(n *memberlist.Node) {
	fmt.Println(c.Name + " " + n.Name + " joined!")
}

func (c client) NotifyLeave(n *memberlist.Node) {

}

func (c client) NotifyUpdate(n *memberlist.Node) {

}

func (c client) NumMembers() int {
	return c.membersList.NumMembers()
}

func (c client) Join(addresses []string) int {
	n, err := c.membersList.Join(addresses)
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}

	return n
}

func (c client) GetMembers() string {
	return "MEMBERS2"
}
