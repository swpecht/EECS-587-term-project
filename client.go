package DUP

import (
	"github.com/hashicorp/memberlist"
	"log"
	"net"
	"strconv"
	"sync"
)

type client struct {
	memberList *memberlist.Memberlist // Underlying memberlist to track membership

	pendingMembersLock sync.Mutex
	pendingMembers     *[]Node // Members that are online, but not active
	ActiveMembersLock  sync.Mutex
	ActiveMembers      map[string]Node // Members that are online and active, mapped by the memberlist.Node.Name
	Name               string          // Unique name of the client
	node               Node            // Used for TCP communications
	MsgChannel         chan Message
	tcpListener        *net.TCPListener
}

func (c client) NotifyJoin(n *memberlist.Node) {
	new_node := Node{
		Name: n.Name,
		Addr: n.Addr,
		Port: int(n.Port) + 100, // Add 100 for the port offset
	}

	if n.Name == c.Name {
		// The initial self notification
		c.ActiveMembersLock.Lock()
		c.ActiveMembers[c.Name] = new_node
		c.ActiveMembersLock.Unlock()
		return
	}

	c.pendingMembersLock.Lock()
	*c.pendingMembers = append(*c.pendingMembers, new_node)
	c.pendingMembersLock.Unlock()
}

func (c client) NotifyLeave(n *memberlist.Node) {

}

func (c client) NotifyUpdate(n *memberlist.Node) {

}

func (c client) NumMembers() int {
	return c.memberList.NumMembers()
}

func (c client) NumActiveMembers() int {
	c.ActiveMembersLock.Lock()
	num := len(c.ActiveMembers)
	c.ActiveMembersLock.Unlock()
	return num
}

// Cause a node to join another memberlist group. This function removes this
// node from the active list. Further more, this should only be called
// when a node is alone in it's undelying memberlist. Therefore, a group
// of nodes cannot merge with another group, but the sub group must all join
// individually. Should this be blocking until the node is made active?
func (c *client) Join(addresses []string) {
	c.memberList.Join(addresses)
	c.updateActiveMemberList([]Node{})
	go c.waitAndActivate()
	return
}

func (c *client) Close() {
	c.tcpListener.Close()
	c.memberList.Shutdown()
}

// Wait until the client is active
func (c *client) WaitActive() {
	for {
		if c.IsActive() == true {
			break
		}
	}
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

func (c *client) updateActiveMemberList(members []Node) {
	log.Println("[DEBUG] Updateing active member list with: " + strconv.Itoa(len(members)))
	c.ActiveMembersLock.Lock()
	c.ActiveMembers = make(map[string]Node) // Reset the active members map

	for i := range members {
		c.ActiveMembers[members[i].Name] = members[i]
		log.Println("[DEBUG] Active member " + strconv.Itoa(i) + " " + members[i].Name)
	}
	c.ActiveMembersLock.Unlock()

}

// Determine if the given client is in the active pool
func (c *client) IsActive() bool {
	c.ActiveMembersLock.Lock()
	_, ok := c.ActiveMembers[c.Name]
	c.ActiveMembersLock.Unlock()
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
