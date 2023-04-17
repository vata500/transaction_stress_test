package main

import (
	"l2_testing_tool/transfertoken"
	"time"
)

func main() {

	checkStartTime := time.Now()
	// sendeth.Start()
	transfertoken.Start(checkStartTime)
	// logging.Start("nitro.log")
}

