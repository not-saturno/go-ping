package network

import (
	"log"
	"net"
)

func LookupDomain(domain string) net.IP {
	ipAddresses, err := net.LookupIP(domain)

	if err != nil {
		log.Fatalf("GO-PING: Error looking up Domain")
	}

	return ipAddresses[0]
}
