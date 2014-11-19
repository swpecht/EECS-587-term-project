package DUP

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
)

// tcpListen listens for and handles incoming connections
func tcpListen(listener *net.TCPListener, channel chan Message) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("[Debug] Closing listener")
			listener.Close()
			break
		}
		go handleConn(conn, channel)
	}
}

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

// handleConn handles a single incoming TCP connection
func handleConn(c *net.TCPConn, channel chan Message) {
	for {
		msg, err := recvMessage(c)
		if err != nil {
			log.Println("[ERROR] Failed to rcvmessage: " + err.Error())
		}
		if err == io.EOF {
			log.Println("[DEBUG] Closing connection.")
			break
		}
		log.Println("[DEBUG] Message recieved ", msg)
		// Quesues messages for processing in the channel
		channel <- msg
	}

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

// Marshal the message and send it over a given TCP connection
func sendMessage(conn *net.TCPConn, msg Message) error {
	// Serialize the message
	msgString, err := msg.Encode()
	if err != nil {
		return err
	}
	log.Println("[DEBUG] Sending message: " + msgString)
	io.Copy(conn, bytes.NewBufferString(msgString))

	return nil
}
