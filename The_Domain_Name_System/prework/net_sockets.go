package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	SERVER_HOSTNAME = "1.1.1.1"
	SERVER_PORT     = 53
	SERVER_PROTOCOL = "udp"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Need to provide 2 arguments! Usage: go run net_sockets.go [domain] [type] ex. google.com A")
		os.Exit(1)
	}

	hostname_to_query := os.Args[1]
	query_type := os.Args[2]

	dns_query := createDNSQuery(0, hostname_to_query, query_type)

	serialized_data := dns_query.serialize()

	fmt.Println("Setting up the connection...")

	conn, err := net.Dial(SERVER_PROTOCOL, SERVER_HOSTNAME+":"+strconv.Itoa(SERVER_PORT))

	if err != nil {
		panic(err)
	}

	fmt.Println("Sending the query...")
	if _, err := conn.Write(serialized_data); err != nil {
		panic(err)
	}

	fmt.Println("Received a response! Reading...")
	read_buffer := make([]byte, 4096)

	for {
		conn.Read(read_buffer)
		dns_response := &DNSMessage{}
		dns_response.deserialize(read_buffer)

		if dns_response.identification != dns_query.identification {
			continue
		}
		fmt.Println("Received a response!")
		dns_response.printDNSMessage()
		break
	}

	conn.Close()

	os.Exit(0)
}
