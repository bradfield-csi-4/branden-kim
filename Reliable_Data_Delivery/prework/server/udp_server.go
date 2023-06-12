package main

import (
	"fmt"
	"net"
	udp_packet "udp_from_scratch/custom_udp_packet"
)

const (
	SERVER_PORT = "9999"
	SERVER_HOST = "127.0.0.1"
)

func main() {

	fmt.Println("Initiating Server Listening on port 9999...")

	conn, err := net.ListenPacket("udp", SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		fmt.Println("Could not open a socket to listen for incoming connections!")
		panic(err)
	}

	defer conn.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT + "...")
	fmt.Println("Waiting for an incoming connection...")
	read_buffer := make([]byte, 4096)
	for {
		fmt.Println("Reading the incoming message...")
		mLen, ip_addr, err := conn.ReadFrom(read_buffer)

		if err != nil {
			fmt.Println("Could not read from the incoming connection!")
			panic(err)
		}

		fmt.Print("Received: ")
		udp_packet.DumpByteSlice(read_buffer[:mLen])
		_, err = conn.WriteTo([]byte("Received message!\n"), ip_addr)
		break
	}
}
