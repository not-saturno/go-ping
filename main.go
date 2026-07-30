package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/not-saturno/go-ping/internal/network"
)

func main() {
	count := flag.Int("c", 1, "Amount of pings to be sent until the end of the program.")

	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalf("usage: go-ping [-c count] <destination>")
	}

	destination := flag.Arg(0)

	ipAddress, err := network.LookupDomain(destination)

	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	pinger, err := network.Listen()

	if err != nil {
		log.Fatalf("GO-PING: Failure starting connection. %v\n", err)
	}

	defer pinger.Close()

	pinger.SetupUDPAddress(ipAddress)

	for range *count {
		nSent, err := pinger.SendMessage()

		if err != nil {
			log.Fatalf("GO-PING: Failure writing message to %s (%s)\n", destination, ipAddress.String())
		}

		fmt.Printf("GO-PING: Sent %d bytes to %s (%s)\n", nSent, destination, ipAddress.String())

		replyMessage, err := pinger.ReadReply()

		if err != nil {
			log.Fatalf("GO-PING: %v", err)
		}

		fmt.Println(replyMessage)
		time.Sleep(1 * time.Second)
	}

}
