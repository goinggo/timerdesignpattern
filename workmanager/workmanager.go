// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This package implements the WorkManager singleton. This manager
	controls the starting, shutdown and processing of work.
*/
package workmanager

import (
	"github.com/goinggo/timerdesignpattern/helper"
	"sync/atomic"
	"time"
)

const (
	TIMER_PERIOD time.Duration = 15 * time.Second // Interval to wake up on
)

// _WorkManager is responsible for starting and shutting down the program
type _WorkManager struct {
	Shutdown        int32
	ShutdownChannel chan string
}

var _This *_WorkManager // Reference to the singleton

// Startup brings the manager to a running state
func Startup() (err error) {

	defer helper.CatchPanic(&err, "main", "workmanager.Startup")

	helper.WriteStdout("main", "workmanager.Startup", "Started")

	// Create the work manager to get the program going
	_This = &_WorkManager{
		Shutdown:        0,
		ShutdownChannel: make(chan string),
	}

	// Start the work timer routine.
	// When this returns the program terminates
	go _This.GoRoutine_WorkTimer()

	helper.WriteStdout("main", "workmanager.Startup", "Completed")

	return err
}

// Shutdown brings down the manager gracefully
func Shutdown() (err error) {

	defer helper.CatchPanic(&err, "main", "workmanager.Shutdown")

	helper.WriteStdout("main", "workmanager.Shutdown", "Started")

	// Shutdown the program
	helper.WriteStdout("main", "workmanager.Shutdown", "Info : Shutting Down")
	atomic.CompareAndSwapInt32(&_This.Shutdown, 0, 1)

	helper.WriteStdout("main", "workmanager.Shutdown", "Info : Shutting Down Work Timer")
	_This.ShutdownChannel <- "Down"
	<-_This.ShutdownChannel

	close(_This.ShutdownChannel)

	helper.WriteStdout("main", "workmanager.Shutdown", "Completed")

	return err
}

// GoRoutine_WorkTimer perform the work on the defined interval
func (this *_WorkManager) GoRoutine_WorkTimer() {

	helper.WriteStdout("WorkTimer", "_WorkManager.GoRoutine_WorkTimer", "Started")

	wait := TIMER_PERIOD

	for {

		helper.WriteStdoutf("WorkTimer", "_WorkManager.GoRoutine_WorkTimer", "Info : Wait To Run : Seconds[%.0f]", wait.Seconds())

		select {

		case <-this.ShutdownChannel:

			helper.WriteStdoutf("WorkTimer", "_WorkManager.GoRoutine_WorkTimer", "Shutting Down")

			this.ShutdownChannel <- "Down"
			return

		case <-time.After(wait):

			helper.WriteStdoutf("WorkTimer", "_WorkManager.GoRoutine_WorkTimer", "Woke Up")
			break
		}

		// Mark the starting time
		startTime := time.Now()

		// Perform the work
		this.PerformTheWork()

		// Mark the ending time
		endTime := time.Now()

		// Caluclate the amount of time to wait to start this again
		duration := endTime.Sub(startTime)
		wait = TIMER_PERIOD - duration
	}
}

// PerformTheWork simulate some silly display work with silly sleep times
func (this *_WorkManager) PerformTheWork() {

	defer helper.CatchPanic(nil, "_WorkManager", "WorkManager.PerformTheWork")

	helper.WriteStdout("WorkTimer", "_WorkManager.GoRoutine_WorkTimer", "Started")

	// Perform work for 10 seconds
	for count := 0; count < 40; count++ {

		if atomic.CompareAndSwapInt32(&_This.Shutdown, 1, 1) == true {

			helper.WriteStdout("WorkTimer", "_WorkManager.GoRoutine_WorkTimer", "Info : Request To Shutdown")
			return
		}

		helper.WriteStdoutf("WorkTimer", "_WorkManager.GoRoutine_WorkTimer", "Processing Images For Station : %d", count)

		time.Sleep(time.Millisecond * 250)
	}

	helper.WriteStdout("WorkTimer", "WorkManager.GoRoutine_WorkTimer", "Completed")
}
