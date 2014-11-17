package DUP

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemberList(t *testing.T) {
	assert := assert.New(t)
	headName := "0.0.0.0:7946"
	// Create some test clients
	factory := ClientFactory{}
	client, _ := factory.NewClient()

	client2, _ := factory.NewClient()
	client3, _ := factory.NewClient()

	client2.Join([]string{headName})
	client3.Join([]string{headName})
	num_clients := client.NumMembers()
	assert.Equal(num_clients, 3, "Incorrect num of initial clients")

	// Test tracking of active nodes
	assert.Equal(1, client.NumActiveMembers(), "bad initial active members")

	num_active := client.UpdateActiveMembers()
	assert.Equal(3, num_active, "invlaid new number of active members.")

	// Test the requirements for all nodes to agree on adding active members
	client4, _ := factory.NewClient()
	client4.Join([]string{headName})

	go client.UpdateActiveMembers()
	assert.Equal(3, client2.NumActiveMembers(), "Allowed actives to join early")

	go client2.UpdateActiveMembers()
	go client3.UpdateActiveMembers()

	assert.Equal(4, client.NumActiveMembers(),
		"New member not allowed to be active")

}

func TestJoining(t *testing.T) {
	t.Errorf("Not implemented.")
}
