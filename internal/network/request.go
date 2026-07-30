package network

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type PacketConn interface {
	ReadFrom(b []byte) (int, net.Addr, error)
	WriteTo(b []byte, addr net.Addr) (int, error)
	Close() error
}

type Pinger struct {
	*icmp.PacketConn
	address *net.UDPAddr
}

func LookupDomain(domain string) (net.IP, error) {
	ipAddresses, err := net.LookupIP(domain)

	if err != nil || len(ipAddresses) == 0 {
		return nil, fmt.Errorf("failed to resolve domain %q: %w", domain, err)
	}

	return ipAddresses[0], nil
}

func NewPinger() {

}

func Listen() (*Pinger, error) {
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")

	if err != nil {
		return nil, err
	}

	return &Pinger{conn, nil}, nil
}

func setupMessage() ([]byte, error) {
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("Echo!"),
		},
	}

	return message.Marshal(nil)
}

func (p *Pinger) SetupUDPAddress(ipAddr net.IP) {
	p.address = &net.UDPAddr{IP: ipAddr, Port: 0}
}

func (p *Pinger) SendMessage() (int, error) {
	message, err := setupMessage()

	if err != nil {
		return 0, err
	}

	return p.WriteTo(message, p.address)
}

func (p *Pinger) ReadReply() (string, error) {
	reply := make([]byte, 1500)

	nRead, peer, err := p.ReadFrom(reply)
	if err != nil {
		log.Fatalf("GO-PING: Read error, %v", err)
	}

	parsedMsg, err := icmp.ParseMessage(1, reply[:nRead])

	if err != nil {
		log.Fatalf("GO-PING Parse error, %v", err)
	}

	switch parsedMsg.Type {
	case ipv4.ICMPTypeEchoReply:
		echoReply, ok := parsedMsg.Body.(*icmp.Echo)
		if !ok {
			return "", fmt.Errorf("invalid ICMP echo reply body format")
		}
		return fmt.Sprintf("Reply from %s: seq=%d bytes=%d", peer, echoReply.Seq, nRead), nil
	default:
		return fmt.Sprintf("Received non-echo ICMP message type: %v", parsedMsg.Type), nil
	}
}
