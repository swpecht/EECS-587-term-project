package DUP

import (
	"github.com/hashicorp/memberlist"
)

type client struct {
	membersList    *memberlist.Memberlist // Underlying memberlist to track membership
	pendingMembers []memberlist.Node      // Members that are online, but not active
	ActiveMembers  []memberlist.Node      // Members that are online and active
	Name           string                 // Unique name of the client
	node           Node                   // Used for TCP communications
}

func (c client) NotifyJoin(n *memberlist.Node) {
	// fmt.Println(c.Name + " " + n.Name + " joined!")
}

func (c client) NotifyLeave(n *memberlist.Node) {

}

func (c client) NotifyUpdate(n *memberlist.Node) {

}

func (c client) NumMembers() int {
	return c.membersList.NumMembers()
}

func (c client) NumActiveMembers() int {
	return len(c.ActiveMembers)
}

// Cause a node to join another memberlist group. This function removes this
// node from the active list. Further more, this should only be called
// when a node is alone in it's undelying memberlist. Therefore, a group
// of nodes cannot merge with another group, but the sub group must all join
// individually.
func (c *client) Join(addresses []string) (int, error) {
	n, err := c.membersList.Join(addresses)
	return n, err
}

// Allows members currently waiting to become active to become active,
// this method blocks and requires that all current active members
// have also called this method.
func (c *client) UpdateActiveMembers() int {
	return c.NumActiveMembers()
}
