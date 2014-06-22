package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func server(startCh chan bool, shutdownCh chan bool) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("(server) Server started. Listening...")
	// Signal server has started
	startCh <- true

	log.Println("(server) Accepting")
	// Wait for a connection to be made
	readConn, acceptErr := ln.Accept()
	if acceptErr != nil {
		log.Fatal(acceptErr)
	}
	log.Println("(server) Got connection")

	log.Println("(server) Reading")
	// Read one line from the connection
	line, readErr := bufio.NewReader(readConn).ReadString('\n')
	if readErr != nil {
		log.Fatal(readErr)
	}

	log.Println("(server) \"" + strings.Trim(line, "\n ") + "\"")

	log.Println("(server) Closing")
	if closeErr := readConn.Close(); closeErr != nil {
		log.Fatal(closeErr)
	}

	log.Println("(server) Bye!")
	// Signal server has finished
	shutdownCh <- true
}

func client() {
	log.Println("(client) Dialing")
	writeConn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("(client) Connected!")

	log.Println("(client) Writing")
	_, writeErr := fmt.Fprintf(writeConn, "Hello!\n")
	if writeErr != nil {
		log.Fatal(writeErr)
	}

	log.Println("(client) Closing")
	if closeErr := writeConn.Close(); closeErr != nil {
		log.Fatal(closeErr)
	}

	log.Println("(client) Bye!")
}

func main() {
	serverStartCh := make(chan bool)
	serverShutdownCh := make(chan bool)

	// Start server in a separate routine
	go server(serverStartCh, serverShutdownCh)
	log.Println("Waiting for server to start")

	// Wait for server to start by blocking on the "start" channel
	<-serverStartCh

	// Launch client in current routine
	client()

	// Wait for server to quit before exiting by blocking on the
	// "shutdown" channel
	<-serverShutdownCh
}
