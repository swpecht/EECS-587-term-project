package DUP

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"
)

// messageType is an integer ID of a type of message that can be received
// on network channels from other members.
type messageType uint8

// The list of available message types.
const (
	activateMsg messageType = iota
	ackMsg
	broadcastMsg
	barrierMsg
)

func CreateBroadcastMsg(stringData []string, floatData []float64) Message {
	msg := Message{
		Type:       broadcastMsg,
		StringData: stringData,
		FloatData:  floatData,
	}

	return msg
}

// Creates an activate message, where the first member of the string array
// contains an array of nodes
func createActivateMsg(activeMembers []Node) (Message, error) {
	nodesBytes, err := json.Marshal(activeMembers)
	if err != nil {
		log.Println("[ERROR] Failed to marshal nodes")
	}
	nodesString := string(nodesBytes)
	msg := Message{
		Type:       activateMsg,
		StringData: []string{nodesString},
	}

	return msg, err
}

func decodeActivateMsg(msg Message) ([]Node, error) {
	var err error
	if msg.Type != activateMsg {
		log.Println("[ERROR] Tried to decodeActivateMsg on non-Activate type message")
		err = errors.New("Failed incorrect message type")
		return make([]Node, 0), err
	}

	nodesString := msg.StringData[0]
	var nodes []Node
	err = json.Unmarshal([]byte(nodesString), &nodes)
	if err != nil {
		log.Println("[ERROR] Failed to unmarshal node list: " + err.Error())
	}

	return nodes, err
}

// Activates all pending members
func (c *client) activatePendingMembers() {
	// Create the appended list of active members
	c.ActiveMembersLock.Lock()
	activeMembers := make([]Node, len(c.ActiveMembers))
	var i int = 0
	for _, value := range c.ActiveMembers {
		activeMembers[i] = value
		i++
	}
	c.ActiveMembersLock.Unlock()

	c.pendingMembersLock.Lock()
	pending_members := *c.pendingMembers
	c.pendingMembersLock.Unlock()
	activeMembers = append(activeMembers, pending_members...)

	msg, _ := createActivateMsg(activeMembers)

	// TODO implement some logic here so everyone does send to the
	// new members
	for i := 0; i < len(pending_members); i++ {
		tcp_conn, _ := c.getTCPConection(pending_members[i])
		tcpAddr := pending_members[i].GetTCPAddr()

		sendMessage(tcp_conn, msg)
		log.Println("[DEBUG] Activate message sent to: ", tcpAddr.String())
	}

	// Update the active members on the local node
	log.Println("[DEBUG] Total active nodes: " + strconv.Itoa(len(activeMembers)))
	c.updateActiveMemberList(activeMembers)
}

// Returns a connection to the specified node
// TODO use a connection pool for speed
func (c *client) getTCPConection(node Node) (*net.TCPConn, error) {
	tcpAddr := node.GetTCPAddr()
	tcp_conn, err := net.DialTCP("tcp", nil, &tcpAddr)
	if err != nil {
		log.Println("[ERROR] Failed to get tcp connection to ", node.GetTCPAddr())
	}

	return tcp_conn, err
}

func (c *client) broadCastMsg(msg Message) {
	c.ActiveMembersLock.Lock()
	log.Println("[DEBUG] Broadcasting message to", len(c.ActiveMembers), "nodes")
	for _, node := range c.ActiveMembers {
		tcpConn, err := c.getTCPConection(node)
		tcpAddr := node.GetTCPAddr()
		if err != nil {
			log.Println("[ERROR] Failed to broadcast message to ", tcpAddr.String())
		}

		err = sendMessage(tcpConn, msg)
		if err != nil {
			log.Println("[ERROR] Failed to broadcast message to ", tcpAddr.String())
		}
	}
	c.ActiveMembersLock.Unlock()

}

// Called to handle the tcp communication of a join.
func (c *client) waitAndActivate() (int, error) {
	activateMessage := <-c.activateChannel
	activeNodes, err := decodeActivateMsg(activateMessage)

	c.updateActiveMemberList(activeNodes)

	log.Println("[DEBUG] Activated node, total active nodes: " + strconv.Itoa(len(activeNodes)))
	return len(activeNodes), err
}
