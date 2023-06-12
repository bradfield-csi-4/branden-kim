package main

import (
	"fmt"
	"net"
	"os"
	udp_packet "udp_from_scratch/custom_udp_packet"
)

const (
	PROXY_IP_ADDRESS = "0.0.0.0"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Need to provide the [MESSAGE] and [PORT] that the UDP proxy server is running on! ex. go run udp_client.go [MESSAGE] [PORT]")
		os.Exit(1)
	}

	udp_client_message := os.Args[1]
	proxy_port := os.Args[2]

	conn, err := net.Dial("udp", PROXY_IP_ADDRESS+":"+proxy_port)

	fmt.Println("Setting up the connection...")
	if err != nil {
		fmt.Println("Could not set up a socket connection to the proxy server!")
		panic(err)
	}

	query := udp_packet.CreateUDPMessage(proxy_port, udp_client_message)

	buffer := query.Serialize()

	udp_packet.DumpByteSlice(buffer)

	fmt.Println("Writing the message over the network...")
	_, err = conn.Write(buffer)

	if err != nil {
		fmt.Println("Could not write to the proxy server!")
		panic(err)
	}

	read_buffer := make([]byte, 4096)
	mLen, err := conn.Read(read_buffer)

	if err != nil {
		fmt.Println("Could not read the response from the proxy server!")
		panic(err)
	}

	fmt.Println("Received: ", string(read_buffer[:mLen]))
	conn.Close()

	os.Exit(0)
}
