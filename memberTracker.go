// Keeps track of membership
// Can be backed by memberlist

package DUP

import (
	"sync"
)

type MemberTracker interface {
	// Get the current memberlist
	GetMemberList() []Node

	// Join a node's list
	Join(target string) error

	// Gracefully leave a member tracker
	Leave() error

	// Start the member list
	Start() error
}

// Implementation of MemberTracker with a syncronized list
// as the backing, useful for testing purposes
// Must be thread safe
type ListMemberTracker struct {
	// This map is shared between all member trackers. The key is the Node
	// name and the value is a bool indicating if the membertracker has joined
	// the other member trackers
	memberList     map[string]bool
	memberListLock *sync.Mutex
	node           Node
}

// Creates a new list member tracker
func New(memberList map[string]bool, lock *sync.Mutex, node Node) ListMemberTracker {
	tracker := ListMemberTracker{
		memberList:     memberList,
		memberListLock: lock,
		node:           node,
	}

	tracker.memberListLock.Lock()
	defer tracker.memberListLock.Unlock()
	// Initialize this tracker to false
	tracker.memberList[tracker.node.Name] = false

	return tracker
}

func (tracker ListMemberTracker) GetMemberList() []string {
	tracker.memberListLock.Lock()
	defer tracker.memberListLock.Unlock()

	list := make([]string, len(tracker.memberList))
	i := 0
	for k, _ := range tracker.memberList {
		list[i] = k
		i++
	}
	return list
}

// Join a member tracker group, for the list backed tracker
// the target string is unused. There is only a single group that
// can be active at a time.
func (tracker ListMemberTracker) Join(target string) error {
	tracker.memberListLock.Lock()
	defer tracker.memberListLock.Unlock()

	tracker.memberList[tracker.node.Name] = true

	return nil
}

// Does nothing for list backed member tracker
func (tracker ListMemberTracker) Start() error {
	return nil
}

func (tracker ListMemberTracker) Leave() error {
	tracker.memberListLock.Lock()
	defer tracker.memberListLock.Unlock()

	tracker.memberList[tracker.node.Name] = false
	return nil
}

type MemberListMemberTracker struct {
}

// func (m MemberListMemberTracker) Start() {

// }
