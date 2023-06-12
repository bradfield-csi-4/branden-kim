package custom_udp_packet

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

const (
	UTF_8 uint16 = 0
)

type CustomUDPApplicationHeader struct {
	sender_port, encoding uint16
	message_length        uint32
}

type CustomUDPApplicationData struct {
	data string
}

type CustomUDPApplicationPacket struct {
	header CustomUDPApplicationHeader
	data   CustomUDPApplicationData
}

func CreateUDPMessage(sender_port string, data string) *CustomUDPApplicationPacket {
	value, _ := strconv.ParseUint(sender_port, 10, 64)
	header := CustomUDPApplicationHeader{
		sender_port:    uint16(value),
		encoding:       UTF_8,
		message_length: uint32(len(data)),
	}

	message := CustomUDPApplicationData{
		data: data,
	}

	return &CustomUDPApplicationPacket{
		header: header,
		data:   message,
	}
}

func DumpByteSlice(b []byte) {
	var a [16]byte
	n := (len(b) + 15) &^ 15
	for i := 0; i < n; i++ {
		if i%16 == 0 {
			fmt.Printf("%4d", i)
		}
		if i%8 == 0 {
			fmt.Print(" ")
		}
		if i < len(b) {
			fmt.Printf(" %02X", b[i])
		} else {
			fmt.Print("   ")
		}
		if i >= len(b) {
			a[i%16] = ' '
		} else if b[i] < 32 || b[i] > 126 {
			a[i%16] = '.'
		} else {
			a[i%16] = b[i]
		}
		if i%16 == 15 {
			fmt.Printf("  %s\n", string(a[:]))
		}
	}
}

func (udp_message *CustomUDPApplicationPacket) Serialize() []byte {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint16(buffer[0:2], udp_message.header.sender_port)
	binary.BigEndian.PutUint16(buffer[2:4], udp_message.header.encoding)
	binary.BigEndian.PutUint32(buffer[4:8], udp_message.header.message_length)

	buffer = append(buffer, []byte(udp_message.data.data)...)

	return buffer
}
