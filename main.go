package main

import (
	"log"
	"runtime"

	"github.com/nmrshll/winter-is-coming-backend/tcpserver"
	"github.com/nmrshll/winter-is-coming-backend/utils/errors"
)

func main() {
	runtime.GOMAXPROCS(1)

	_, stoppedChan, err := tcpserver.StartTCPServer(7777)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed starting tcp server"))
	}
	<-stoppedChan // keep server open for the lifetime of the program
}
