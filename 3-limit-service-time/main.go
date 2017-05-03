//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import "time"

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	resultChan := make(chan bool)

	before := time.Now().Unix()

	go func() {
		process()
		resultChan <- true
	}()

	var result bool

	if u.IsPremium {
		result = <-resultChan
	} else {
		select {
		case result = <-resultChan:
		case <-time.After(time.Second * time.Duration(10-u.TimeUsed)):
			result = false
		}
	}

	// Update the consumed time
	after := time.Now().Unix()
	requestTime := after - before
	u.TimeUsed = u.TimeUsed + requestTime

	return result
}

func main() {
	RunMockServer()
}
