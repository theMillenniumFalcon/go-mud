package main

import (
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) error {
	log.Println("I got a connection")
	return nil
}

func startServer() error {
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
			if err := handleConnection(conn); err != nil {
				log.Println("Error handling connection", err)
				return
			}
		}()
	}
}

func main() {
	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}
