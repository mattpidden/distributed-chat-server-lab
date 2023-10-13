package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// Deal with an error event.
	fmt.Println("Server has encountered an error: ", err.Error())
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		conn, _ := ln.Accept()
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)
	for {
		msg, _ := reader.ReadString('\n')
		if msg != "\n" {
			newMsg := Message{
				sender:  clientid, // Assign a value to sender
				message: msg,      // Assign a value to message
			}
			//fmt.Fprintln(client, "Server Recived.")
			msgs <- newMsg
		}
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//Create a Listener for TCP connections on the port given above.
	ln, _ := net.Listen("tcp", *portPtr)
	fmt.Println("Running Server on Port", *portPtr)


	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)
	numberClients := 0
	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			// - Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client
			numberClients += 1
			fmt.Println("New Client Connected, ID:", numberClients)
			clients[numberClients] = conn
			go handleClient(conn, numberClients, msgs)
		case msg := <-msgs:
			// - Deal with a new message
			fmt.Println("MESSAGE RECIVED USER", msg.sender, "::", msg.message)
			// Send the message to all clients that aren't the sender
			for clientId, connection := range clients {
				if clientId != msg.sender {
					fmt.Fprintln(connection, "User", msg.sender, "::", msg.message)
				}
			}
		}
	}
}
