package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type clientType map[string]bool

var clients = clientType{}


func (c clientType) sharePeers(conn *net.UDPConn, toAddr string) {

	var peers []string
	for k, _ := range c {
		if k != toAddr {
			peers = append(peers, k)
		}
	}
	if len(peers) == 0 {
		return
	}

	msg := strings.Join(peers, ",")
	r, _ := net.ResolveUDPAddr("udp", toAddr)
	conn.WriteTo([]byte("++"+msg), r)

	fmt.Printf("> sending ++%s to %s\n", msg, toAddr)

}

func (c clientType) shareWithPeers(conn *net.UDPConn, addr string, extra string) {
	for k := range c {
		if k != addr {
			r, _ := net.ResolveUDPAddr("udp", k)
			conn.WriteTo([]byte(extra+addr), r)
			fmt.Printf("> sending %s%s to %s\n", extra, addr, k)
		}
	}

}

func Server() {

	// set default port if none provided
	localAddress := ":9999"
	if len(os.Args) > 2 {
		localAddress = os.Args[2]
	}

	// create a udp server
	addr, _ := net.ResolveUDPAddr("udp", localAddress)
	conn, _ := net.ListenUDP("udp", addr)

	fmt.Printf("Listening %s...\n", localAddress)

	for {

		// read udp message
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		addrStr := addr.String()

		
		switch string(buf[:n]) {
		// receiived regestration message
		case "!reg":
			fmt.Printf(">%s %s from %s%s\n", Green, buf[:n], addr, Reset)
			// add address to know peers list
			clients[addrStr] = true
			// share all peers with newly connected peer
			clients.sharePeers(conn, addrStr)
			// share newly connected peer with all other peers
			clients.shareWithPeers(conn, addrStr, "++")
		// received de-registration message
		case "!dereg":
			fmt.Printf(">%s %s from %s%s\n", Red, buf[:n], addr, Reset)
			// delete from know peers list
			delete(clients, addrStr)
			// tell peers to remove address from their peer list
			clients.shareWithPeers(conn, addrStr, "--")
		}

	}

}
