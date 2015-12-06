package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
)

func handleConnection(server *Server, conn net.Conn) {
	server.openConnection(conn)
}

func main() {
	flag.Parse()
	runtime.SetBlockProfileRate(1)
	serverDbPath := filepath.Join(persistDir, "dispatchd-server.db")
	msgDbPath := filepath.Join(persistDir, "messages.db")
	var server = NewServer(serverDbPath, msgDbPath)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", amqpPort))
	if err != nil {
		fmt.Printf("Error!\n")
		os.Exit(1)
	}
	fmt.Printf("Listening on port %d\n", amqpPort)
	go func() {
		fmt.Printf("Go perf handlers on port %d\n", debugPort)
		log.Println(http.ListenAndServe(fmt.Sprintf("localhost:%d", debugPort), nil))
	}()
	go func() {
		startAdminServer(server, adminPort)
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection!\n")
			os.Exit(1)
		}
		go handleConnection(server, conn)
	}
}

var amqpPort int
var debugPort int
var adminPort int
var persistDir string

func init() {
	flag.IntVar(&amqpPort, "amqp-port", 1111, "Port for amqp protocol messages")
	flag.IntVar(&debugPort, "debug-port", 6060, "Port for the golang debug handlers")
	flag.IntVar(&adminPort, "admin-port", 8080, "Port for admin server")
	flag.StringVar(&persistDir, "persist-dir", "/data/dispatchd/", "Directory for the server and message database files")
}
