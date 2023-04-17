package main

import (
	"flag"
	"l2_testing_tool/sendeth"
	"l2_testing_tool/transfertoken"
	"time"
)

func main() {
	transferPtr := flag.Bool("transfer", false, "Run transfer method")
	sendPtr := flag.Bool("send", false, "Run send method")
	loggingPtr := flag.Bool("logging", false, "Run logging method")
	flag.Parse()

	if *transferPtr {
		checkStartTime := time.Now()
		transfertoken.Start(checkStartTime)
	}
	if *sendPtr {
		sendeth.Start()
	}
	if *loggingPtr {
		sendeth.Start()
	}



	// logging.Start("nitro.log")
}

