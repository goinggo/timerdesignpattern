// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This program provides an sample to learn how to implement a timer
	routine and graceful shutdown pattern.

	Ardan Studios
	12973 SW 112 ST, Suite 153
	Miami, FL 33186
	bill@ardanstudios.com

	http://www.goinggo.net/2013/09/timer-routines-and-graceful-shutdowns.html
*/
package main

import (
	"bufio"
	"github.com/goinggo/timerdesignpattern/helper"
	"github.com/goinggo/timerdesignpattern/workmanager"
	"os"
)

// main is the starting point of the program
func main() {

	helper.WriteStdout("main", "main", "Starting Program")

	workmanager.Startup()

	// Hit enter to terminate the program gracefully
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	workmanager.Shutdown()

	helper.WriteStdout("main", "main", "Program Complete")
}
