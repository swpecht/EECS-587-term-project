package DUP

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
)

type Messenger interface {
	Send(msg Message) error
	Recv(channel chan Message) error
	resolve(addr string) (interface{}, error)
}

type ChannelMessenger struct {
	ResolverMap map[string]chan Message // Used by the resolver function to send messages
	Incoming    chan Message
}

func (messenger *ChannelMessenger) Send(msg Message) error {
	var resolved interface{}
	resolved, err := messenger.resolve(msg.Target)
	if err != nil {
		return err
	}
	targetChan, ok := resolved.(chan Message)
	if !ok || targetChan == nil {
		errorMsg := "Failed to convert the resolved value to a channel"
		err = errors.New(errorMsg)
		log.Println("[ERROR]", errorMsg, resolved)
		return err
	}

	log.Println("[DEBUG] Sending message", msg, "over", targetChan)
	targetChan <- msg
	return nil
}

func (messenger *ChannelMessenger) Recv(channel chan Message) error {
	log.Println("[DEBUG] Waiting on receive on", messenger.Incoming)
	incomingMessage := <-messenger.Incoming
	log.Println("[DEBUG] Message received", incomingMessage)
	channel <- incomingMessage
	return nil
}

func (messenger *ChannelMessenger) resolve(addr string) (interface{}, error) {
	channel, ok := messenger.ResolverMap[addr]
	var err error
	if !ok {
		log.Println("[ERROR] Failed to resolve", addr)
		err = errors.New("Address not found!")
	}
	log.Println("[DEBUG] Resolved", addr, "to", channel)
	return channel, err
}

// tcpListen listens for and handles incoming connections
func tcpListen(listener *net.TCPListener, channel chan Message) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("[DEBUG] Closing listener", listener.Addr, err.Error())
			// listener.Close()
			return
		}
		go handleConn(conn, channel)
	}
}

type Message struct {
	Type       messageType
	Target     string
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
	msg, err := recvMessage(c)
	if err != nil {
		log.Println("[ERROR] Failed to rcvmessage: " + err.Error())
	}
	if err == io.EOF {
		log.Println("[DEBUG] Closing connection.")
		c.Close()
		return
	}
	// log.Println("[DEBUGMessage recieved ", msg)
	// Quesues messages for processing in the channel
	channel <- msg

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
	//This is probably timing out, may beed to use a thread ppol
	//_, err = conn.Write([]byte(msgString))
	return err
}
