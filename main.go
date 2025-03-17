// tunnel_client.go
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
)

var (
	serverAddr string
	targetAddr string
)

func init() {
	flag.StringVar(&serverAddr, "server", "", "Tunnel server address (<host>:<port>)")
	flag.StringVar(&targetAddr, "target", "", "Target address (<host>:<port>)")
}

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", ":4000") // Port where the tunnel client listens
	if err != nil {
		panic(err)
	}
	fmt.Println("Tunnel client listening on port 4000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleRequest(conn)
	}
}

func handleRequest(clientConn net.Conn) {
	fmt.Println("Received connection from server")

	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		fmt.Println("Error connecting to target:", err)
		clientConn.Close()
		return
	}

	fmt.Println("Connected to target")
	go copyIO(clientConn, targetConn)
	go copyIO(targetConn, clientConn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(dest, src)
}
