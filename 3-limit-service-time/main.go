package main

import (
	"sync"
	"time"
)

const SECONDS_PER_USER = 10

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	lock      sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {

	processDone := make(chan struct{})

	go func() {
		process()
		processDone <- struct{}{}

	}()

	for {
		select {
		case <-time.After(time.Second):
			u.lock.Lock()
			if !u.IsPremium && u.TimeUsed >= 10 {
				u.lock.Unlock()
				return false
			}
			u.TimeUsed += 1
			u.lock.Unlock()
		case <-processDone:
			return true
		}
	}
}

func main() {
	RunMockServer()
}
