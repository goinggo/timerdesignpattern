// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package helper provides helper functions for logging and catching panics.
package helper

import (
	"fmt"
	"runtime"
	"time"
)

// WriteStdout is used to write message directly stdout.
func WriteStdout(goRoutine string, functionName string, message string) {
	fmt.Printf("%s : %s : %s : %s\n", time.Now().Format("2006-01-02T15:04:05.000"), goRoutine, functionName, message)
}

// WriteStdoutf is used to write a formatted message directly stdout.
func WriteStdoutf(goRoutine string, functionName string, format string, a ...interface{}) {
	WriteStdout(goRoutine, functionName, fmt.Sprintf(format, a...))
}

// CatchPanic is used to catch and display panics.
func CatchPanic(err *error, goRoutine string, function string) {
	if r := recover(); r != nil {
		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		WriteStdoutf(goRoutine, function, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}
