package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	log.Println("Launching GoKV server...")

	ln, errListen := net.Listen("tcp", ":3334")

	if errListen != nil {
		log.Fatal(errListen.Error())
		os.Exit(1)
	}

	defer ln.Close()

	log.Println("Listening on port 3334")

	for {
		conn, errAccept := ln.Accept()

		if errAccept != nil {
			log.Fatal(errAccept.Error())
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		newmessage := strings.ToUpper(message)
		conn.Write([]byte(newmessage + "\n"))
	}
}
