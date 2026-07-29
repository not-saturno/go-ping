package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalf("usage: go-ping [-c count] <destination>")
	}

	destination := flag.Arg(0)

	ipAddr, err := net.LookupIP(destination)

	if err != nil {
		log.Fatalf("GO-PING Name or service not found.")
	}

	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatalf("GO-PING Listen error: %v", err)
	}

	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}

	bytes, err := msg.Marshal(nil)
	dstAddr := &net.UDPAddr{IP: ipAddr[0], Port: 0}

	nSent, err := conn.WriteTo(bytes, dstAddr)
	fmt.Printf("GO-PING Sent %d bytes to %s (%s)\n", nSent, destination, ipAddr[0].String())

	reply := make([]byte, 1500)
	nRead, peer, err := conn.ReadFrom(reply)
	if err != nil {
		log.Fatalf("GO-PING Read error, %v", err)
	}

	parsedMsg, err := icmp.ParseMessage(1, reply[:nRead])
	if err != nil {
		log.Fatalf("GO-PING Parse error, %v", err)
	}

	switch parsedMsg.Type {
	case ipv4.ICMPTypeEchoReply:
		echoReply, ok := parsedMsg.Body.(*icmp.Echo)
		if !ok {
			log.Fatalf("GO-PING invalid ICMP echo body format")
		}
		fmt.Printf("Reply from %s (%s): seq=%d bytes=%d", destination,
			peer, echoReply.Seq, nRead)

	default:
		fmt.Printf("GO-PING Got non-echo ICMP message type: %v\n", parsedMsg.Type)
	}

}
