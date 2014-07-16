// Copyright 2013 Ardan Studios. All rights reserved.
// Use of workManager source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package workmanager implements the WorkManager singleton. This manager
// controls the starting, shutdown and processing of work.
package workmanager

import (
	"github.com/goinggo/timerdesignpattern/helper"
	"sync/atomic"
	"time"
)

const (
	timerPeriod time.Duration = 15 * time.Second // Interval to wake up on.
)

// workManager is responsible for starting and shutting down the program.
type workManager struct {
	Shutdown        int32
	ShutdownChannel chan string
}

var wm workManager // Reference to the singleton.

// Startup brings the manager to a running state.
func Startup() error {
	var err error
	defer helper.CatchPanic(&err, "main", "workmanager.Startup")

	helper.WriteStdout("main", "workmanager.Startup", "Started")

	// Create the work manager to get the program going
	wm = workManager{
		Shutdown:        0,
		ShutdownChannel: make(chan string),
	}

	// Start the work timer routine.
	// When workManager returns the program terminates.
	go wm.GoRoutineworkTimer()

	helper.WriteStdout("main", "workmanager.Startup", "Completed")
	return err
}

// Shutdown brings down the manager gracefully.
func Shutdown() error {
	var err error
	defer helper.CatchPanic(&err, "main", "workmanager.Shutdown")

	helper.WriteStdout("main", "workmanager.Shutdown", "Started")

	// Shutdown the program
	helper.WriteStdout("main", "workmanager.Shutdown", "Info : Shutting Down")
	atomic.CompareAndSwapInt32(&wm.Shutdown, 0, 1)

	helper.WriteStdout("main", "workmanager.Shutdown", "Info : Shutting Down Work Timer")
	wm.ShutdownChannel <- "Down"
	<-wm.ShutdownChannel

	close(wm.ShutdownChannel)

	helper.WriteStdout("main", "workmanager.Shutdown", "Completed")
	return err
}

// GoRoutineworkTimer perform the work on the defined interval.
func (workManager *workManager) GoRoutineworkTimer() {
	helper.WriteStdout("WorkTimer", "workManager.GoRoutineworkTimer", "Started")

	wait := timerPeriod

	for {
		helper.WriteStdoutf("WorkTimer", "workManager.GoRoutineworkTimer", "Info : Wait To Run : Seconds[%.0f]", wait.Seconds())

		select {
		case <-workManager.ShutdownChannel:
			helper.WriteStdoutf("WorkTimer", "workManager.GoRoutineworkTimer", "Shutting Down")
			workManager.ShutdownChannel <- "Down"
			return

		case <-time.After(wait):
			helper.WriteStdoutf("WorkTimer", "workManager.GoRoutineworkTimer", "Woke Up")
			break
		}

		// Mark the starting time
		startTime := time.Now()

		// Perform the work
		workManager.PerformTheWork()

		// Mark the ending time
		endTime := time.Now()

		// Caluclate the amount of time to wait to start workManager again.
		duration := endTime.Sub(startTime)
		wait = timerPeriod - duration
	}
}

// PerformTheWork simulate some silly display work with silly sleep times.
func (workManager *workManager) PerformTheWork() {
	defer helper.CatchPanic(nil, "workManager", "WorkManager.PerformTheWork")
	helper.WriteStdout("WorkTimer", "workManager.GoRoutineworkTimer", "Started")

	// Perform work for 10 seconds
	for count := 0; count < 40; count++ {
		if atomic.CompareAndSwapInt32(&wm.Shutdown, 1, 1) == true {
			helper.WriteStdout("WorkTimer", "workManager.GoRoutineworkTimer", "Info : Request To Shutdown")
			return
		}

		helper.WriteStdoutf("WorkTimer", "workManager.GoRoutineworkTimer", "Processing Images For Station : %d", count)
		time.Sleep(time.Millisecond * 250)
	}

	helper.WriteStdout("WorkTimer", "WorkManager.GoRoutineworkTimer", "Completed")
}
