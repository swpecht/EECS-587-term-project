package DUP

import (
	"github.com/hashicorp/memberlist"
	"net"
)

type client struct {
	memberList     *memberlist.Memberlist // Underlying memberlist to track membership
	pendingMembers *[]Node                // Members that are online, but not active
	ActiveMembers  map[string]Node        // Members that are online and active, mapped by the memberlist.Node.Name
	Name           string                 // Unique name of the client
	node           Node                   // Used for TCP communications
	MsgChannel     chan Message
	tcpListener    *net.TCPListener
}

func (c client) NotifyJoin(n *memberlist.Node) {
	if n.Name == c.Name {
		// The initial self notification
		c.ActiveMembers[n.Name] = c.node
		return
	}
	new_node := Node{
		Addr: n.Addr,
		Port: int(n.Port) + 100, // Add 100 for the port offset
	}
	*c.pendingMembers = append(*c.pendingMembers, new_node)
}

func (c client) NotifyLeave(n *memberlist.Node) {

}

func (c client) NotifyUpdate(n *memberlist.Node) {

}

func (c client) NumMembers() int {
	return c.memberList.NumMembers()
}

func (c client) NumActiveMembers() int {
	return len(c.ActiveMembers)
}

// Cause a node to join another memberlist group. This function removes this
// node from the active list. Further more, this should only be called
// when a node is alone in it's undelying memberlist. Therefore, a group
// of nodes cannot merge with another group, but the sub group must all join
// individually. Should this be blocking until the node is made active?
func (c *client) Join(addresses []string) (int, error) {
	c.memberList.Join(addresses)
	c.ActiveMembers = make(map[string]Node) // Reset the active members map
	n, err := c.waitAndActivate()
	return n, err
}

// Allows members currently waiting to become active to become active,
// this method blocks and requires that all current active members
// have also called this method.
func (c *client) UpdateActiveMembers() int {
	// Need to ensure all active members have decided to do this
	c.Barrier()
	c.activatePendingMembers()
	// Need to send go ahead message to new members to be made active
	return c.NumActiveMembers()
}

// Determine if the given client is in the active pool
func (c *client) IsActive() bool {
	_, ok := c.ActiveMembers[c.Name]
	return ok
}

// Barrier that blacks for all active nodes
func (c *client) Barrier() {
	if !c.IsActive() {
		panic("Client is not active!")
	}
	if len(c.ActiveMembers) == 1 {
		// This is the only member so can return instantly
		return
	}
	return
}

// Send message to all nodes
func (c *client) Broadcast(stringData []string, floatData []float64) {

}
