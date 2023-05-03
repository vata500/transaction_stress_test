package main

import (
	"flag"
	"fmt"
	"l2_testing_tool/logging"
	"l2_testing_tool/transfertoken"
	"time"
)

func main() {
	transferPtr := flag.Bool("transfer", false, "Run transfer method")
	loggingPtr := flag.Bool("logging", false, "Run logging method")
	flag.Parse()

	if *transferPtr {
		checkStartTime := time.Now()
		fmt.Printf("%s",checkStartTime)
		transfertoken.Start(checkStartTime)
	}
	if *loggingPtr {
		logging.Start()
	}
}

