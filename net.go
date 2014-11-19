package DUP

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
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
	ackMsg
	broadcastMsg
)

// Messages are sent with the first byte being the message type
// and a string body that represents JSON that can be decerialized
// into the appropriate message type.
type Message struct {
	Type       messageType
	StringData []string
	FloatData  []float64
}

// Encodes a messafe for sending over a tcp connection. Format is:
// {len in}\n{msgbody}
func (msg Message) Encode() (outputMsg string, err error) {
	msgBody, err := json.Marshal(msg)
	if err != nil {
		log.Println("[ERROR] Failed to encode message: " + err.Error())
	}
	outputMsg += string(msgBody) + string('\n')
	return
}

func Decode(b []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(b, &msg)
	if err != nil {
		log.Println("[ERROR] Failed to unmarshal message")
	}

	return msg, err
}

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

// tcpListen listens for and handles incoming connections
func (c *client) tcpListen(channel chan Message) {
	for {
		conn, err := c.tcpListener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn, channel)
	}
}

// Activates all pending members
func (c *client) activatePendingMembers() {
	// Create the appended list of active members
	c.ActiveMembersLock.RLock()
	activeMembers := make([]Node, len(c.ActiveMembers))
	var i int = 0
	for _, value := range c.ActiveMembers {
		activeMembers[i] = value
		i++
	}
	c.ActiveMembersLock.RUnlock()

	c.pendingMembersLock.RLock()
	pending_members := *c.pendingMembers
	c.pendingMembersLock.RUnlock()
	activeMembers = append(activeMembers, pending_members...)

	msg, _ := createActivateMsg(activeMembers)

	// TODO implement some logic here so everyone does send to the
	// new members
	for i := 0; i < len(pending_members); i++ {
		tcpAddr := pending_members[i].GetTCPAddr()
		tcp_conn, _ := net.DialTCP("tcp", nil, &tcpAddr)
		sendMessage(tcp_conn, msg)
		log.Println("[DEBUG] Activate message sent to: " + tcpAddr.String())
	}

	// Update the active members on the local node
	log.Println("[DEBUG] Total active nodes: " + strconv.Itoa(len(activeMembers)))
	c.updateActiveMemberList(activeMembers)
}

// handleConn handles a single incoming TCP connection
func handleConn(c *net.TCPConn, channel chan Message) {
	msg, err := recvMessage(c)
	if err != nil {
		log.Println("[ERROR] Failed to rcvmessage: " + err.Error())
	}
	log.Println("[DEBUG] Message recieved.")
	channel <- msg
}

// Marshal the message and send it over a given TCP connection
func sendMessage(conn *net.TCPConn, msg Message) error {
	// Serialize the message
	msgString, err := msg.Encode()
	if err != nil {
		return err
	}
	log.Println("[DEBUG] Serialized Message: " + msgString)
	io.Copy(conn, bytes.NewBufferString(msgString))

	return nil
}

// Receive a message over a tcp connections, and unmarshal it from JSON
func recvMessage(conn *net.TCPConn) (Message, error) {
	reader := bufio.NewReader(conn)
	b, err := reader.ReadBytes('\n')
	if err != nil {
		log.Println("[ERROR] Failed to read message")
		return Message{}, err
	}
	msg, err := Decode(b)
	return msg, err
}

// Called to handle the tcp communication of a join.
func (c *client) waitAndActivate() (int, error) {
	activateMessage := <-c.MsgChannel
	activeNodes, err := decodeActivateMsg(activateMessage)

	c.updateActiveMemberList(activeNodes)

	log.Println("[DEBUG] Activated node, total active nodes: " + strconv.Itoa(len(activeNodes)))
	return len(activeNodes), err
}
