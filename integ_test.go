package DUP

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMemberList(t *testing.T) {
	assert := assert.New(t)
	headName := "0.0.0.0:7946"
	// Create some test clients
	factory := ClientFactory{}
	client, _ := factory.NewClient()

	client2, _ := factory.NewClient()
	client3, _ := factory.NewClient()

	defer client.Close()
	defer client2.Close()
	defer client3.Close()

	client2.Join([]string{headName})
	client3.Join([]string{headName})
	num_clients := client.NumMembers()
	assert.Equal(num_clients, 3, "Incorrect num of initial clients")

	// Test tracking of active nodes
	assert.Equal(1, client.NumActiveMembers(), "bad initial active members")
	assert.Equal(0, client2.NumActiveMembers(), "Not purging active after join")

	// Test tracking of pending nodes
	assert.Equal(2, len(*client.pendingMembers), "Not tracking pending members")

	num_active := client.UpdateActiveMembers()
	assert.Equal(3, num_active, "invlaid new number of active members.")
	client2.WaitActive()
	client3.WaitActive()
	assert.Equal(3, client2.NumActiveMembers(), "invlaid new number of active members.")
	assert.Equal(3, client3.NumActiveMembers(), "invlaid new number of active members.")

	assert.True(client2.IsActive())
	assert.True(client3.IsActive())

	// Can only join when not already part of a group
	// Test the requirements for all nodes to agree on adding active members
	// Tests the barrier functionality
	client4, _ := factory.NewClient()
	defer client3.Close()
	client4.Join([]string{headName})

	go client.UpdateActiveMembers()
	assert.Equal(3, client2.NumActiveMembers(), "Allowed actives to join early")

	go client2.UpdateActiveMembers()
	go client3.UpdateActiveMembers()

	time.Sleep(5000 * time.Millisecond)
	//client4.WaitActive()
	assert.Equal(4, client.NumActiveMembers(),
		"New member not allowed to be active")

	t.Errorf("Not implemented.")

}

func TestActiveStatus(t *testing.T) {
	assert := assert.New(t)
	headName := "0.0.0.0:7946"
	// Create some test clients
	factory := ClientFactory{}
	client, err := factory.NewClient()
	if err != nil {
		t.Errorf("Failed to create client: " + err.Error())
	}
	client2, err := factory.NewClient()
	if err != nil {
		t.Errorf("Failed to create client: " + err.Error())
	}

	defer client.Close()
	defer client2.Close()

	client2.Join([]string{headName})
	assert.False(client2.IsActive())
	client.UpdateActiveMembers()
	client2.WaitActive()

	assert.True(client2.IsActive())
}
