// Keeps track of membership
// Can be backed by memberlist

package DUP

import ()

type MemberTracker interface {
	// Get the current memberlist
	GetMemberList() []Node

	// Join a node's list
	Join(target string) error

	// Gracefully leave a member tracker
	Leave() error

	// Start the member list
	Start() error
}

// Implementation of MemberTracker with a syncronized list
// as the backing, useful for testing purposes
type ListMemberTracker struct {
}

type MemberListMemberTracker struct {
}

// func (c client) NotifyJoin(n *memberlist.Node) {
// 	new_node := Node{
// 		Name: n.Name,
// 		Addr: n.Addr,
// 		Port: int(n.Port) + 100, // Add 100 for the port offset
// 	}

// 	if n.Name == c.Name {
// 		// The initial self notification
// 		c.ActiveMembersLock.Lock()
// 		c.ActiveMembers[c.Name] = new_node
// 		c.ActiveMembersLock.Unlock()
// 		return
// 	}

// 	c.pendingMembersLock.Lock()
// 	*c.pendingMembers = append(*c.pendingMembers, new_node)
// 	c.pendingMembersLock.Unlock()
// }

// func (c client) NotifyLeave(n *memberlist.Node) {

// }

// func (c client) NotifyUpdate(n *memberlist.Node) {

// }

// func (m MemberListMemberTracker) Start() {
// 	var config *memberlist.Config = memberlist.DefaultLocalConfig()
// 	config.BindPort = c.node.Port - 100 // off set for tcp
// 	config.Name = c.Name
// 	config.AdvertisePort = c.node.Port - 100 // off set for tcp
// 	config.Events = c

// 	list, err := memberlist.Create(config)
// 	if err != nil {
// 		log.Println("[ERROR] Failed to create member list for", c.Name, "Error: ", err.Error())
// 	}
// 	log.Println("[DEBUG] Started memberlist for", c.Name)
// 	c.memberList = list
// }
