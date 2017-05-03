//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import "os"
import "os/signal"
import "syscall"

func main() {
	// Subscibe on SIGINTs
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a process
	proc := MockProcess{}

	go proc.Run()

	// Try to shut down gracefully
	<-sigs
	println("\n Gracefully shutting down...")
	gracefulShutDown := make(chan bool)
	go func() {
		proc.Stop()
		gracefulShutDown <- true
	}()

	// Wait for a process shutdown or another SIGINT
	select {
	case <-gracefulShutDown:
		println("\n Gracefully shut down!")
	case <-sigs:
		println("\n Killed on demand")
	}
}
