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
    startCh <- true

    log.Println("(server) Accepting")
    readConn, acceptErr := ln.Accept()
    if acceptErr != nil {
        log.Fatal(acceptErr)
    }
    log.Println("(server) Got connection")

    log.Println("(server) Reading")
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
    go server(serverStartCh, serverShutdownCh)
    log.Println("Waiting for server to start")
    <-serverStartCh
    client()
    <-serverShutdownCh
}
