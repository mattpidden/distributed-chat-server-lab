package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func read(conn net.Conn, inboundChannel chan string) {
	//In a continuous loop, read a message from the server and display it.
	reader := bufio.NewReader(conn)
	for {
		msg, _ := reader.ReadString('\n')
		if msg != "\n" {
			inboundChannel <- msg
		}
	}
}

func write(conn net.Conn, outboundChannel chan string) {
	//Continually get input from the user and send messages to the server.
	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter Your Message: ")
		msg, _ := stdin.ReadString('\n')
		outboundChannel <- msg

	}


}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()
	//Try to connect to the server
	conn, _ := net.Dial("tcp", *addrPtr)
	fmt.Println("Connected to Server on", *addrPtr)
	inboundChannel := make(chan string)
	outboundChannel := make(chan string)
	go read(conn, inboundChannel)
	go write(conn, outboundChannel)
	for {
		select {
			case inbound := <-inboundChannel:
				fmt.Println("\n", inbound)
			case outbound := <-outboundChannel:
				fmt.Fprintln(conn, outbound)
		}
	}
}
