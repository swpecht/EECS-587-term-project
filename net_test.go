package DUP

import (
	"log"
	"net"
	"testing"
	"time"
)

func TestSendRecvMessage(t *testing.T) {

	listenAddr, err := net.ResolveTCPAddr("tcp", "localhost:5000")
	if err != nil {
		t.Errorf("Couldn't resolve hostname")
	}

	log.Println("[INFO] Creating listener...")
	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		t.Errorf("Couldn't create listener")
	}

	log.Println("[INFO] Starting listener...")
	c := make(chan Message)
	go HandleIncomeing(listener, c)

	remoteConn, _ := net.DialTCP("tcp", nil, listenAddr)

	msg := Message{
		Type:       activateMsg,
		StringData: []string{"Hello World"},
	}

	err = sendMessage(remoteConn, msg)
	if err != nil {
		t.Errorf(err.Error())
	}

	go func() {
		time.Sleep(1 * time.Second)
		c <- Message{StringData: []string{"ERROR"}}
	}()
	msgReceived := <-c
	log.Println("[DEBUG] Recieved: " + msgReceived.StringData[0])
	if msgReceived.StringData[0] == "ERROR" {
		t.Error("Timed out!")
	}

}

func TestEncodeDecodeMessage(t *testing.T) {
	weights := []float64{1.0, 2.5, 3}
	strings := []string{"hello", "world"}
	msg := CreateBroadcastMsg(strings, weights)
	msgString, err := msg.Encode()
	log.Println("[DEBUG] encoded message: " + msgString)
	if err != nil {
		t.Error("Encoding message failed")
	}

	decodedMsg, err := Decode([]byte(msgString))
	if err != nil {
		t.Error("Decoding message failed")
	}
	if decodedMsg.FloatData[0] != msg.FloatData[0] ||
		decodedMsg.FloatData[1] != msg.FloatData[1] ||
		decodedMsg.FloatData[2] != msg.FloatData[2] {
		t.Error("Decoded floats not the same")
	}

	if decodedMsg.StringData[0] != msg.StringData[0] ||
		decodedMsg.StringData[1] != msg.StringData[1] {
		t.Error("Decoded strings not the same")
	}

	if decodedMsg.Type != broadcastMsg {
		t.Error("Incorrect message type")
	}

}

func HandleIncomeing(l *net.TCPListener, c chan Message) {
	for {
		// Wait for a connection.
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleConn(conn, c)
	}
}
