package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	header "traceroute/header"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Need to specify an argument! Pass in a hostname. ex. go run traceroute.go [HOSTNAME]")
		os.Exit(1)
	}

	hostname := os.Args[1]

	fmt.Println("Looking up IP Address...")
	addr, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println("Could not look up IP address of the hostname")
		panic(err)
	}

	fmt.Println("Opening Socket...")
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)

	ipv4_addr := net.IP(addr[0])
	p := pkt()

	sockaddr := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{8, 8, 8, 8},
		// Addr: ipv4_addr[0:4],
	}

	if err != nil {
		fmt.Println("Could not open a socket!")
		syscall.Close(fd)
		panic(err)
	}

	fmt.Println(ipv4_addr)
	// err = syscall.Connect(fd, &sockaddr)

	// if err != nil {
	// 	fmt.Printf("Could not bind socket!")
	// 	panic(err)
	// }

	fmt.Println("Sending the ICMP packet...")
	err = syscall.Sendto(fd, p, 0, &sockaddr)
	if err != nil {
		log.Fatal("Sendto:", err)
		syscall.Close(fd)
	}

	fmt.Println("Reading the ICMP packet...")
	read_buffer := make([]byte, 1024)
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	for {

		numRead, err := f.Read(read_buffer)

		if err != nil {
			fmt.Println("Could not receive the packet!")
			syscall.Close(fd)
			panic(err)
		}
		fmt.Printf("% X\n", read_buffer[:numRead])
		break
	}

	syscall.Close(fd)
	os.Exit(0)
}

func pkt() []byte {
	h := header.Header{
		Version:  4,
		Len:      20,
		TotalLen: 20 + 10, // 20 bytes for IP, 10 for ICMP
		TTL:      64,
		Protocol: 1, // ICMP
		Dst:      net.IPv4(8, 8, 8, 8),
		// ID, Src and Checksum will be set for us by the kernel
	}

	icmp := []byte{
		8, // type: echo request
		0, // code: not used by echo request
		0, // checksum (16 bit), we fill in below
		0,
		0, // identifier (16 bit). zero allowed.
		0,
		0, // sequence number (16 bit). zero allowed.
		0,
		0xC0, // Optional data. ping puts time packet sent here
		0xDE,
	}
	cs := csum(icmp)
	icmp[2] = byte(cs)
	icmp[3] = byte(cs >> 8)

	_, err := h.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	return append(icmp)
}

func csum(b []byte) uint16 {
	var s uint32
	for i := 0; i < len(b); i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	// add back the carry
	s = s>>16 + s&0xffff
	s = s + s>>16
	return uint16(^s)
}
