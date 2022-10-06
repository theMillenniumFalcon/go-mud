package main

import (
	"fmt"
	"log"
	"net"
)

type MessageEvent struct {
	msg string
}

type UserJoinedEvents struct {
	user *User
}

type User struct {
	name    string
	session *Session
}

type Session struct {
	conn *net.Conn
}

type World struct {
	users []*User
}

func handleConnection(conn net.Conn, eventChannel chan interface{}) error {
	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return err
		}
		if n == 0 {
			log.Println("Zero bytes received, closing connection")
			break
		}
		msg := string(buf[0 : n-2])
		log.Println("Received Message:", msg)

		e := &MessageEvent{msg}
		eventChannel <- e

		resp := fmt.Sprintf("You said, \"%s\"\r\n", msg)
		n, err = conn.Write([]byte(resp))
		if err != nil {
			return err
		}
		if n == 0 {
			log.Println("Zero bytes received, closing connection")
			break
		}
	}
	return nil
}

func startServer(eventChannel <-chan interface{}) error {
	fmt.Println("Starting server")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}
		go func() {
			if err := handleConnection(conn, eventChannel); err != nil {
				log.Println("Error handling connection", err)
				return
			}
		}()
	}
}

func startGameLoop(ch <-chan interface{}) {
	// w := &World{}
	for event := range ch {
		fmt.Println("Recieved event", event)
	}
}

func main() {
	ch := make(chan interface{})

	go startGameLoop(ch)

	err := startServer(ch)
	if err != nil {
		log.Fatal(err)
	}
}
