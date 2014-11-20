package DUP

import (
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
)

// Mock messenger for client tests. Collects all sent messages into a channel
// Does no resolving. Has a capacity for 50 sent messages
type MockMessenger struct {
	sentMessages chan Message
}

func NewMockMessenger() (Messenger, chan Message) {
	msgChan := make(chan Message, 50)
	return MockMessenger{
		sentMessages: msgChan,
	}, msgChan
}

func (messenger MockMessenger) Send(msg Message) error {
	messenger.sentMessages <- msg
	return nil
}

func (messenger MockMessenger) Recv(channel chan Message) error {
	return nil
}

// Does nothing
func (messenger MockMessenger) resolve(addr string) (interface{}, error) {
	return nil, nil
}

var nodeNumLock sync.Mutex
var nodeNum int = 0

func GetNode(t *testing.T) Node {
	nodeNumLock.Lock()
	addr := net.ParseIP("0.0.0.0")

	node := Node{
		Name: string(nodeNum),
		Port: nodeNum,
		Addr: addr,
	}
	nodeNum++
	nodeNumLock.Unlock()
	return node
}

func GetClient_DataOnly(t *testing.T) *client {
	c := new(client)
	f := ClientFactory{}
	f.initializeData(c)

	c.ActiveMembers[c.node.Name] = c.node
	return c
}

func GetActivateMessage(t *testing.T, nodes []Node) Message {
	msg, err := createActivateMsg(nodes)
	if err != nil {
		t.Errorf("Failed to create activate message")
	}

	return msg
}

func GetBarrierMessage(t *testing.T, source string) Message {
	msg := createBarrierMsg(source)
	return msg
}

func TestClient_IsActive(t *testing.T) {
	assert := assert.New(t)
	c := GetClient_DataOnly(t)

	assert.True(c.IsActive())

	c.updateActiveMemberList([]Node{})
	assert.False(c.IsActive())

}

func TestClient_Barrier_Blocking(t *testing.T) {
	t.Errorf("Not Implemented")
}

func TestClient_HandleActivate(t *testing.T) {
	assert := assert.New(t)
	c := GetClient_DataOnly(t)
	emptyActivate := GetActivateMessage(t, []Node{})
	c.handleActivateMessage(emptyActivate)
	assert.Equal(0, c.NumActiveMembers())

	node := GetNode(t)
	sameActivate := GetActivateMessage(t, []Node{node, node})
	c.handleActivateMessage(sameActivate)
	assert.Equal(1, c.NumActiveMembers())

	twoActivate := GetActivateMessage(t, []Node{GetNode(t), GetNode(t)})
	c.handleActivateMessage(twoActivate)
	assert.Equal(2, c.NumActiveMembers())

}

func TestClient_Close(t *testing.T) {
	t.Errorf("Not Implemented")
}
