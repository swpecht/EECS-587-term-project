package DUP

import (
	"bytes"
	"encoding/json"
	"io"
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
	ack
)

// Messages are sent with the first byte being the message type
// and a string body that represents JSON that can be decerialized
// into the appropriate message type.
type Message struct {
	Type messageType
	Data string
}

// Encodes a messafe for sending over a tcp connection. Format is:
// {len in}\n{msgbody}
func (msg Message) Enconde() (outputMsg string, err error) {
	msgBody, err := json.Marshal(msg)

	msgLen := len(msgBody)
	outputMsg = strconv.Itoa(msgLen) + string('\n')
	outputMsg += string(msgBody) + string('\n')

	return
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
	// TODO implement locks for this
	pending_members := *c.pendingMembers
	activeMembers := make([]Node, len(c.ActiveMembers)+len(pending_members))

	activeMembers = append(activeMembers, pending_members...)

	for _, value := range c.ActiveMembers {
		activeMembers = append(activeMembers, value)
	}

	for i := 0; i < len(pending_members); i++ {
		tcpAddr := pending_members[i].GetTCPAddr()
		tcp_conn, _ := net.DialTCP("tcp", nil, &tcpAddr)
		c.sendActivateMessage(tcp_conn, activeMembers)
	}
}

// Send a message over a given TCP connection
func sendMessage(conn *net.TCPConn, msg Message) error {
	// Serialize the message
	msgString, err := msg.Enconde()
	if err != nil {
		return err
	}
	log.Println("[DEBUG] Serialized Message: " + msgString)
	io.Copy(conn, bytes.NewBufferString(msgString))

	return nil
}

// Sends an activate message with the specified active nodes over the given
// tcp connection.
func (c *client) sendActivateMessage(conn *net.TCPConn, activeNodes []Node) error {
	activateMessage := ActivateMsg{ActiveNodes: activeNodes}
	b, err := json.Marshal(activateMessage)
	if err != nil {
		return err
	}
	msg := Message{Type: activateMsg, Data: string(b)}
	return sendMessage(conn, msg)
}

// Called to handle the tcp communication of a join.
func (c *client) waitAndActivate() (int, error) {
	return 0, nil
}
