package DUP

import (
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
)

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
	f.startMessageHandling(c)

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

// http://stackoverflow.com/questions/19167970/mock-functions-in-golang
// Good thoughts on how to mock out the messageing
// https://eclipse.org/paho/clients/golang/
// http://golang.org/pkg/net/textproto/
// Good messaging library
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
