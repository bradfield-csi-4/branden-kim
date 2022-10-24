package main

import (
	"fmt"
	"os"
	"syscall"
)

var (
	SOCKET_PORT = 53
	SOCKET_ADDR = [4]byte{8, 8, 8, 8}
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Need to provide 2 arguments! Usage: go run net_sockets.go [domain] [type] ex. google.com A")
		os.Exit(1)
	}

	hostname_to_query := os.Args[1]
	query_type := os.Args[2]

	dns_query := createDNSQuery(0, hostname_to_query, query_type)

	serialized_data := dns_query.serialize()

	sock_struct := syscall.SockaddrInet4{
		Port: SOCKET_PORT,
		Addr: SOCKET_ADDR,
	}

	fmt.Println("Setting up the connection...")

	fd, sock_err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)

	if sock_err != nil {
		syscall.Close(fd)
		panic(sock_err)
	}

	connect_err := syscall.Connect(fd, &sock_struct)

	if connect_err != nil {
		if connect_err == syscall.ECONNREFUSED {
			fmt.Println("* Connection Refused")
			syscall.Close(fd)
		}
		panic(connect_err)
	}

	fmt.Printf("Server: Bound to addr: %d, port: %d\n", sock_struct.Addr, sock_struct.Port)

	_, send_err := syscall.Write(fd, serialized_data)

	if send_err != nil {
		syscall.Close(fd)
		panic(send_err)
	}

	read_buffer := make([]byte, 4096)

	for {
		_, from, recv_err := syscall.Recvfrom(fd, read_buffer, 0)

		if recv_err != nil {
			syscall.Close(fd)
			panic(recv_err)
		}

		fromip4, ok := from.(*syscall.SockaddrInet4)
		if !ok {
			continue
		}

		if fromip4.Addr != SOCKET_ADDR || fromip4.Port != SOCKET_PORT {
			continue
		}

		dns_response := &DNSMessage{}
		dns_response.deserialize(read_buffer)

		if dns_response.identification != dns_query.identification {
			continue
		}

		fmt.Println("Received a response!")
		dns_response.printDNSMessage()
		break
	}

	syscall.Close(fd)
	os.Exit(0)
}
