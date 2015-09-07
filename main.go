package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

type Operation interface {
	Execute(map[string]string) string
}

type GetOperation struct {
	key string
}

func (op *GetOperation) Execute(store map[string]string) string {
	value, _ := store[op.key]
	return value
}

type SetOperation struct {
	key, value string
}

func (op *SetOperation) Execute(store map[string]string) string {
	store[op.key] = op.value
	return ""
}

type message struct {
	data []string
	conn net.Conn
}

var channel = make(chan message)
var store map[string]string = make(map[string]string)

func main() {
	log.Println("Launching GoKV server...")

	ln, errListen := net.Listen("tcp", ":3334")

	if errListen != nil {
		log.Fatal(errListen.Error())
		os.Exit(1)
	}

	defer ln.Close()

	log.Println("Listening on port 3334...")

	go resolveCommand()

	for {
		conn, errAccept := ln.Accept()

		if errAccept != nil {
			log.Println(errAccept.Error())
		}

		go handleConnection(conn)
	}
}

func resolveCommand() {
	for {
		select {
		case message := <-channel:
			var operation Operation

			if message.data[0] == "set" {
				operation = &SetOperation{key: message.data[1], value: message.data[2]}
			}

			if message.data[0] == "get" {
				operation = &GetOperation{key: message.data[1]}
			}

			result := operation.Execute(store)
			message.conn.Write([]byte(result + "\n"))
		}
	}
}

func handleConnection(conn net.Conn) {
	for {
		m, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			log.Println(err.Error())
		}

		m = strings.Trim(m, " ")

		if m == "" {
			conn.Write([]byte("Unknown command\n"))
			continue
		}

		channel <- message{data: strings.Split(m, " "), conn: conn}
	}
}
