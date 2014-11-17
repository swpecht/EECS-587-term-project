package DUP

import (
	"encoding/json"
	"net"
)

// messageType is an integer ID of a type of message that can be received
// on network channels from other members.
type messageType uint8

// The list of available message types.
const (
	activateMsg messageType = iota
	ack
)

// Messages are sent with the first byte being the message type
// and a string body that represents JSON that can be decerialized
// into the appropriate message type.
type Message struct {
	Type messageType
	Data interface{}
}

type ActivateMsg struct {
	ActiveNodes []Node
}

// tcpListen listens for and handles incoming connections
func (c *client) tcpListen() {
	for {
		conn, _ := c.tcpListener.AcceptTCP()
		go c.handleConn(conn)
	}
}

// handleConn handles a single incoming TCP connection
func (c *client) handleConn(conn *net.TCPConn) {
	defer conn.Close()

}

// Activates all pending members
func (c *client) activatePendingMembers() {
	pending_members := *c.pendingMembers
	for i := 0; i < len(pending_members); i++ {
		tcpAddr := pending_members[i].GetTCPAddr()
		tcp_conn, _ := net.DialTCP("tcp", nil, &tcpAddr)

	}
}

// Send a message over a given TCP connection
func sendMessage(conn *net.TCPConn, msg Message) error {
	// Serialize the message
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	conn.Write(b)
	return err
}

// Sends an activate message with the specified active nodes over the given
// tcp connection.
func (c *client) sendActivateMessage(conn *net.TCPConn, activeNodes []Node) error {
	activateMessage := ActivateMsg{ActiveNodes: activeNodes}
	msg := Message{Type: activateMsg, Body: activateMessage}
	return sendMessage(conn, msg)
}

// Called to handle the tcp communication of a join.
func (c *client) waitAndActivate() (int, error) {
	return 0, nil
}
