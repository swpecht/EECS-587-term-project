package DUP

import (
	"github.com/hashicorp/memberlist"
)

type client struct {
	memberList     *memberlist.Memberlist // Underlying memberlist to track membership
	pendingMembers []*memberlist.Node     // Members that are online, but not active
	ActiveMembers  map[string]Node        // Members that are online and active, mapped by the memberlist.Node.Name
	Name           string                 // Unique name of the client
	node           Node                   // Used for TCP communications
}

func (c client) NotifyJoin(n *memberlist.Node) {
	// fmt.Println(c.Name + " " + n.Name + " joined!")
	// Maybe send an event rather than this, as the proper information
	// is not present here.
	//
	// The active members could be stored as a map, based on the memberlist.Node,
	// That way, the membership could be updated fairly easily, based on the information here.
	// Then, when a node joins, it chould scatter it's connection info to everyone.
	// Then all nodes would have the updated active memberlist, and would be able to continue after the
	// UpdateActiveMembers call.

	if n.Name == c.Name {
		// The initial self notification
		c.ActiveMembers[n.Name] = c.node
		return
	}

	c.pendingMembers = append(c.pendingMembers, n)
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
	n, err := c.memberList.Join(addresses)
	c.ActiveMembers = make(map[string]Node) // Reset the active members map
	return n, err
}

// Allows members currently waiting to become active to become active,
// this method blocks and requires that all current active members
// have also called this method.
func (c *client) UpdateActiveMembers() int {
	return c.NumActiveMembers()
}
