package DUP

import (
	"github.com/hashicorp/memberlist"
	"log"
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

	msgIncoming chan Message // Main channel on which to receive messages
	msgOutgoing chan Message // channel on which to send messages

	barrierChannel  chan string // The channel that handles barrier message, will be the name of the node that sent the barrier
	activateChannel chan Message
	closeChannel    chan bool // channel to stop processing messages
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
	return
}

// Handles different types of messages
func (c *client) startMessageHandling() {

	for {
		select {
		case <-c.closeChannel:
			break // done processing
		case msg := <-c.msgIncoming:
			log.Println("[DEBUG]", c.node.Name, " Message received", msg)
			switch msg.Type {
			case activateMsg:
				c.activateChannel <- msg
				break
			case barrierMsg:
				c.barrierChannel <- msg.StringData[0] // Pass on the name, will be handled on the calling thread
				break
			default:
				// log.Println("[ERROR] Unknown message type")
			}
		}

	}
}

func (c *client) startActivateHandling() {
	for {
		select {
		case <-c.closeChannel:
			break // done processing
		case msg := <-c.activateChannel:
			c.handleActivateMessage(msg)
		}

	}
}

func (c *client) handleActivateMessage(msg Message) {

	activeNodes, err := decodeActivateMsg(msg)
	if err != nil {
		log.Println("[ERROR] Received malformed activate message")
		return
	}
	c.updateActiveMemberList(activeNodes)
	log.Println("[DEBUG] Activated node, total active nodes: " + strconv.Itoa(len(activeNodes)))
}

func (c *client) Close() {
	// c.tcpListener.Close()
	c.memberList.Shutdown()
	c.closeChannel <- true
	// Not totally sure how closing channels works TODO
	close(c.msgIncoming) // Let all blocked processes finish
	close(c.barrierChannel)
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
	if c.NumActiveMembers() == 1 {
		// This is the only member so can return instantly
		return
	}

	// Need to broadcast the barrier message
	msg := createBarrierMsg(c.Name)

	log.Println("[DEBUG] Broadcasting barrier from", c.Name)
	c.broadCastMsg(msg)

	// Wait for each node to respond
	//Get messages from channel

	responded := make(map[string]bool)
PollingLoop:
	for {
		select {
		case name := <-c.barrierChannel:
			responded[name] = true
			log.Println("[DEBUG]", c.Name, "Received barrier from", name, len(responded), "of", c.NumActiveMembers())
		default:
			if len(responded) == c.NumActiveMembers() {
				log.Println("[DEBUG] Barrier completed by", c.Name)
				break PollingLoop // everyone is at the barrier
			}
		}
	}
}

// Send message to all nodes
// TODO implement a tree rather than naive send to all
func (c *client) Broadcast(stringData []string, floatData []float64) {
	msg := CreateBroadcastMsg(stringData, floatData)
	c.broadCastMsg(msg)
}
