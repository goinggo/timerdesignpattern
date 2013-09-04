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
