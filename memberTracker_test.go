package DUP

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func GetListMemberTrackers(num int) []ListMemberTracker {
	return nil
}

func TestMemberTracker_List(t *testing.T) {
	assert := assert.New(t)
	num := 5
	trackers := GetListMemberTrackers(num)
	assert.Equal(num, len(trackers))

	// Test Starting and initialization
	for _, v := range trackers {
		v.Start()
		assert.Equal(1, len(v.GetMemberList()))
	}

	// Test Joining

	// Test Leaving

	t.Error("Not implemented.")
}
