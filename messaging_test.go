package DUP

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func GetChannelMessengers(num int) []ChannelMessenger {
	messengers := make([]ChannelMessenger, num)

	// Generate resolver map
	resolverMap := make(map[string]chan Message)
	for i := 0; i < num; i++ {
		name := "Messenger" + strconv.Itoa(i)
		channel := make(chan Message)
		resolverMap[name] = channel

		messengers[i] = ChannelMessenger{}
		messengers[i].Incoming = channel
	}

	// Update the messengers resolver map
	for i := 0; i < num; i++ {
		messengers[i].ResolverMap = make(map[string]chan Message)
		for k, v := range resolverMap {
			messengers[i].ResolverMap[k] = v
		}
	}

	return messengers
}

func TestMessaging_ChannelMesseger(t *testing.T) {
	assert := assert.New(t)
	messengers := GetChannelMessengers(2)
	messenger0 := messengers[0]
	messenger1 := messengers[1]

	msgTo1 := Message{
		Target:     "Messenger1",
		StringData: []string{"Message to 1"},
	}
	recvrChannel := make(chan Message)

	timer := time.AfterFunc(500*time.Millisecond, func() {
		panic("Hung sending message!")
	})
	defer timer.Stop()

	go messenger1.Recv(recvrChannel)
	err := messenger0.Send(msgTo1)
	if err != nil {
		t.Errorf("Failed to send message", err.Error())
	}
	msgRecvd := <-recvrChannel

	assert.Equal(msgTo1, msgRecvd)

	go messenger0.Send(msgTo1)
	go messenger1.Recv(recvrChannel)
	msgRecvd = <-recvrChannel
	assert.Equal(msgTo1, msgRecvd)

	// Check invalid send
	invalidMessage := Message{
		Target: "Fake Messenger",
	}
	err = messenger0.Send(invalidMessage)
	if err == nil {
		t.Errorf("Failed to handle invalid address on send")
	}
}
