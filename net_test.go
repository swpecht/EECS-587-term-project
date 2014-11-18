package DUP

import (
	"bufio"
	"log"
	"net"
	"testing"
)

func TestSendMessage(t *testing.T) {

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
	c := make(chan string)
	go HandleIncomeing(listener, c)

	remoteConn, _ := net.DialTCP("tcp", nil, listenAddr)

	msg := Message{
		Type: activateMsg,
		Data: []byte("Hello World"),
	}

	err = sendMessage(remoteConn, msg)
	if err != nil {
		t.Errorf(err.Error())
	}

	dataReceived := <-c

	log.Println("[DEBUG] " + dataReceived)
	t.Error("Not implemented")

}

func TestUnwrapMessage(t *testing.T) {
	t.Errorf("Not implemented")
}

func HandleIncomeing(l *net.TCPListener, c chan string) {
	for {
		// Wait for a connection.
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn, channel chan string) {
			reader := bufio.NewReader(conn)
			data, err := reader.ReadString('\n')
			if err != nil {
				log.Println("[ERROR] Failed read")
			}
			stringData := data
			log.Println("[DEBUG] Recieved: " + stringData)
			channel <- stringData
		}(conn, c)
	}
}
