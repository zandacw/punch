package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var peersList = make(map[string]string)

func Client() {

	conn, err := register()
	if err != nil {
		panic(err)
	}

	listen(conn)

}

func register() (*net.UDPConn, error) {

	// handle interupt signals
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// set default port if none provided
	signalAddress := os.Args[2]
	localAddress := ":8888"
	if len(os.Args) > 3 {
		localAddress = os.Args[3]
	}

	// create rendezvous udp address
	remote, _ := net.ResolveUDPAddr("udp", signalAddress)
	// create local udp server
	local, _ := net.ResolveUDPAddr("udp", localAddress)
	conn, _ := net.ListenUDP("udp", local)

	fmt.Printf("> %ssignal server %s%s\n", Cyan, remote, Reset)

	// send registration message to the rendezvous server
	data := []byte("!reg")
	_, err := conn.WriteTo(data, remote)
	if err != nil {
		return nil, err
	}

	// handle interuupts in goroutine
	go func() {
		<-sigchnl
		data := []byte("!dereg")
		conn.WriteTo(data, remote)
		fmt.Printf("> dereg %s\n", remote)
		os.Exit(1)
	}()

	fmt.Printf("> listening %s\n", local)

	return conn, nil

}


func periodicPunch(conn *net.UDPConn, delay int) {
	for {
		time.Sleep(time.Second * time.Duration(delay))
		for peer := range peersList {
			go punch(conn, peer, false)
		}
	}
}

func listen(conn *net.UDPConn) {

	// run chat interface in goroutine
	go chat(conn)

	// periodically punch udp holes for all know peers
	go periodicPunch(conn, 10)

	for {

		// read udp message
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		buf = buf[:n]

		// handle new peer joined
		if n > 2 && string(buf)[:2] == "++" {

			fmt.Printf("> received %s%s%s from signal server\n", Green, buf, Reset)

			peers := strings.Split(strings.ReplaceAll(string(buf), "++", ""), ",")
			for _, peer := range peers {
				peersList[peer] = ""
				punch(conn, peer, true)
			}

			continue

		// handle peer left
		} else if n > 2 && string(buf)[:2] == "--" {

			fmt.Printf("> received %s%s%s from signal server\n", Red, buf, Reset)

			peers := strings.Split(strings.ReplaceAll(string(buf), "--", ""), ",")
			for _, peer := range peers {
				delete(peersList, peer)
			}

			continue

		// ignore data as it was a punch from peer
		} else if string(buf) == "." {

			continue

		// set name of peer
		} else if strings.Contains(string(buf), "#name ") {
			
			msgSplit := strings.Split(string(buf), "#name ")
			name := msgSplit[1]
			peersList[addr.String()] = name

		// handle chat message
		} else {

			name, ok := peersList[addr.String()]
			if ok && name != "" {
				fmt.Printf("%s[%s] %s%s\n", Blue, name, buf, Reset)
			} else {
				fmt.Printf("%s[%s] %s%s\n", Blue, addr.IP, buf, Reset)
			}

		}

	}

}

func punch(conn *net.UDPConn, peer string, output bool) {
	peerAddr := resolvePeerAddr(peer)
	if output {
		fmt.Printf(">%s punch %s%s\n", Yellow, peerAddr, Reset)
	}
	conn.WriteTo([]byte("."), peerAddr)
}

func chat(conn *net.UDPConn) {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		for peer := range peersList {
			peerAddr := resolvePeerAddr(peer)
			conn.WriteTo(scanner.Bytes(), peerAddr)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}

func resolvePeerAddr(peer string) *net.UDPAddr {
	peerAddr, _ := net.ResolveUDPAddr("udp", peer)
	if strings.Contains(peer, "127.0.0.1") {
		peerSpl := strings.Split(peer, ":")
		peerAddr, _ = net.ResolveUDPAddr("udp", "54.91.215.145:"+peerSpl[1])
	}
	return peerAddr
}
