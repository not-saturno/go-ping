<div align="center">

# Go-Ping

Go-Ping is a Go-based networking tool that recreates the core functionality of the standard ping utility. Built to explore low-level network programming, raw socket handling, and ICMP protocol mechanics in Go, it offers a clean, reliable way to probe host reachability and measure network performance.

</div>

## Prerequisites

* **Go 1.22+** installed on your system.

## Installation

Clone the repository and build the executable binary using Go:

git clone https://github.com/not-saturno/go-ping.git
cd go-ping
go build -o go-ping main.go

Alternatively, you can install it directly into your $GOPATH/bin:

go install github.com/not-saturno/go-ping@latest

## Usage

Because `go-ping` utilizes raw sockets to construct and listen for ICMP packets, administrative or root privileges (`sudo`) are required to run the binary.

sudo ./go-ping [-c count] <destination>

### Flags

* `-c`: Amount of pings to be sent until the end of the program (default: 1).

### Examples

Ping a domain once:
sudo ./go-ping google.com

Send 4 pings to an IP address:
sudo ./go-ping -c 4 8.8.8.8