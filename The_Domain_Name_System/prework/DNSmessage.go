package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
)

var DNSQueryTypeMap = map[string]uint16{
	"A":     1,
	"NS":    2,
	"CNAME": 5,
	"SOA":   6,
	"AAAA":  28,
}

type QuestionRecord struct {
	name  string
	qtype uint16
	class uint16
}

type ResourceRecord struct {
	name     string
	rrtype   uint16
	class    uint16
	ttl      uint32
	rdlength uint16
	rdata    []byte
}

type DNSMessage struct {
	identification, qr, opcode, aa, tc, rd, ra, zero, rcode uint16
	number_of_questions, number_of_answers_rrs              uint16
	number_of_authority_rrs, number_of_additional_rrs       uint16
	questions                                               []QuestionRecord
	resource_records                                        []ResourceRecord
}

func createDNSQuery(opcode uint16, hostname_to_query string, query_type string) *DNSMessage {
	questions_list := make([]QuestionRecord, 1)
	questions_list[0] = QuestionRecord{
		name:  hostname_to_query,
		qtype: DNSQueryTypeMap[query_type],
		class: 1,
	}
	resource_records_list := make([]ResourceRecord, 0)
	return &DNSMessage{
		identification:           uint16(rand.Uint32()),
		qr:                       0,
		opcode:                   opcode,
		aa:                       0,
		tc:                       0,
		rd:                       1,
		ra:                       1,
		zero:                     0,
		rcode:                    0,
		number_of_questions:      1,
		number_of_answers_rrs:    0,
		number_of_authority_rrs:  0,
		number_of_additional_rrs: 0,
		questions:                questions_list,
		resource_records:         resource_records_list,
	}
}

func (dns *DNSMessage) serialize() []byte {
	buffer := make([]byte, 256)

	// Encoding the Identification number
	buffer[0] = byte(dns.identification >> 8)
	buffer[1] = byte(dns.identification)

	var lower_flags byte = 0
	var higher_flags byte = 0

	lower_flags = lower_flags | byte(dns.qr)
	lower_flags = lower_flags | (byte(dns.opcode) << 1)
	lower_flags = lower_flags | (byte(dns.aa) << 5)
	lower_flags = lower_flags | (byte(dns.tc) << 6)
	lower_flags = lower_flags | (byte(dns.rd) << 7)

	higher_flags = higher_flags | byte(dns.ra)
	higher_flags = higher_flags | (byte(dns.zero) << 1)
	higher_flags = higher_flags | (byte(dns.rcode) << 4)

	// Encoding the flags
	buffer[2] = higher_flags
	buffer[3] = lower_flags

	// Encoding the number of questions
	buffer[4] = byte(dns.number_of_questions >> 8)
	buffer[5] = byte(dns.number_of_questions)

	// Encoding the number of answers
	buffer[6] = byte(dns.number_of_answers_rrs >> 8)
	buffer[7] = byte(dns.number_of_answers_rrs)

	// Encoding the number of authorities
	buffer[8] = byte(dns.number_of_authority_rrs >> 8)
	buffer[9] = byte(dns.number_of_authority_rrs)

	// Encoding the number of additional
	buffer[10] = byte(dns.number_of_additional_rrs >> 8)
	buffer[11] = byte(dns.number_of_additional_rrs)

	buffer_index := 12
	// Encoding the Questions
	for _, question := range dns.questions {
		split_name := strings.Split(question.name, ".")
		for _, string_to_encode := range split_name {
			length_of_string := len(string_to_encode)
			buffer[buffer_index] = byte(length_of_string)
			buffer_index++
			for _, char_in_string := range string_to_encode {
				buffer[buffer_index] = byte(char_in_string)
				buffer_index++
			}
		}
		buffer[buffer_index] = byte(0)
		buffer_index++
		buffer[buffer_index] = byte(question.qtype >> 8)
		buffer_index++
		buffer[buffer_index] = byte(question.qtype)
		buffer_index++
		buffer[buffer_index] = byte(question.class >> 8)
		buffer_index++
		buffer[buffer_index] = byte(question.class)
		buffer_index++
	}

	// Don't need to serialize the Resource Records since queries don't use them

	return buffer
}

func (dns *DNSMessage) deserialize(buffer []byte) *DNSMessage {
	// Decoding the identification number
	id_number := binary.BigEndian.Uint16(buffer[0:2])

	// Decoding the flags
	flags := binary.BigEndian.Uint16(buffer[2:4])
	qr_flag := flags & 0x0001
	op_flag := flags & 0x001e
	aa_flag := flags & 0x0020
	tc_flag := flags & 0x0040
	rd_flag := flags & 0x0080
	ra_flag := flags & 0x0100
	zeros_flag := flags & 0x0e00
	rcode_flag := flags & 0xf000

	// Decoding the Question Count
	number_of_questions := binary.BigEndian.Uint16(buffer[4:6])

	// Decoding the Answer Count
	number_of_answers := binary.BigEndian.Uint16(buffer[6:8])

	// Decoding the Authorities Count
	number_of_authorities := binary.BigEndian.Uint16(buffer[8:10])

	// Decoding the Additional Count
	number_of_additional := binary.BigEndian.Uint16(buffer[10:12])

	// Decoding the Questions
	var buffer_index uint16 = 12
	questions_list := make([]QuestionRecord, 0)
	for count := 0; count < int(number_of_questions); count++ {
		domain_name := ""
		for {
			subdomain_length := buffer[buffer_index]
			buffer_index++
			if subdomain_length == 0 {
				break
			}
			domain_name = domain_name + "." + string(buffer[buffer_index:buffer_index+uint16(subdomain_length)])
			buffer_index += uint16(subdomain_length)
		}
		question_type := binary.BigEndian.Uint16(buffer[buffer_index : buffer_index+2])
		buffer_index += 2
		question_class := binary.BigEndian.Uint16(buffer[buffer_index : buffer_index+2])
		buffer_index += 2
		questions_list = append(questions_list, QuestionRecord{name: domain_name, qtype: question_type, class: question_class})
	}

	// Decoding the Answer Resource Records
	resource_records_list := make([]ResourceRecord, 0)
	for count := 0; count < int(number_of_answers); count++ {
		domain_name := ""
		for {
			subdomain_length := buffer[buffer_index]
			// This case means that the response is using message compression and that the values lie somewhere else
			if (subdomain_length | 0x3f) == 0xff {
				// Parse out the byte address value to jump to
				pointer_to_address := binary.BigEndian.Uint16(buffer[buffer_index:buffer_index+2]) & 0x3fff
				parseResourceRecordName(buffer, pointer_to_address, &domain_name)
				buffer_index += 2
				break
			} else {
				buffer_index++
				if subdomain_length == 0 {
					break
				}
				domain_name = domain_name + "." + string(buffer[buffer_index:buffer_index+uint16(subdomain_length)])
				buffer_index += uint16(subdomain_length)
			}
		}
		answer_type := binary.BigEndian.Uint16(buffer[buffer_index : buffer_index+2])
		buffer_index += 2
		answer_class := binary.BigEndian.Uint16(buffer[buffer_index : buffer_index+2])
		buffer_index += 2
		answer_ttl := binary.BigEndian.Uint32(buffer[buffer_index : buffer_index+4])
		buffer_index += 4
		answer_rdlength := binary.BigEndian.Uint16(buffer[buffer_index : buffer_index+2])
		buffer_index += 2
		answer_rddata := buffer[buffer_index : buffer_index+uint16(answer_rdlength)]
		buffer_index += uint16(answer_rdlength)
		resource_records_list = append(resource_records_list, ResourceRecord{name: domain_name, rrtype: answer_type, class: answer_class, ttl: answer_ttl, rdlength: answer_rdlength, rdata: answer_rddata})
	}

	dns.identification = id_number
	dns.qr = qr_flag
	dns.opcode = op_flag
	dns.aa = aa_flag
	dns.tc = tc_flag
	dns.rd = rd_flag
	dns.ra = ra_flag
	dns.zero = zeros_flag
	dns.rcode = rcode_flag
	dns.number_of_questions = number_of_questions
	dns.number_of_answers_rrs = number_of_answers
	dns.number_of_authority_rrs = number_of_authorities
	dns.number_of_additional_rrs = number_of_additional
	dns.questions = questions_list
	dns.resource_records = resource_records_list

	return dns
}

func parseResourceRecordName(buffer []byte, pointer_to_address uint16, result_string *string) {
	for {
		subdomain_length := buffer[pointer_to_address]
		pointer_to_address++
		if subdomain_length == 0 {
			break
		}
		*result_string = *result_string + "." + string(buffer[pointer_to_address:pointer_to_address+uint16(subdomain_length)])
		pointer_to_address += uint16(subdomain_length)
	}
}

func (dns *DNSMessage) printDNSMessage() {
	fmt.Println("====================== DNS Header ========================")
	fmt.Println("HEADER INFO")
	fmt.Printf("Identification Number: %d\n", dns.identification)
	fmt.Printf("QR Flag: %d\n", dns.qr)
	fmt.Printf("Opcode Flag: %d\n", dns.opcode)
	fmt.Printf("AA Flag: %d\n", dns.aa)
	fmt.Printf("TC Flag: %d\n", dns.tc)
	fmt.Printf("RD Flag: %d\n", dns.rd)
	fmt.Printf("RA Flag: %d\n", dns.ra)
	fmt.Printf("RCODE Flag: %d\n", dns.rcode)
	fmt.Println("NUMBER OF QUESTIONS / ANSWERS")
	fmt.Printf("Number of Questions: %d\n", dns.number_of_questions)
	fmt.Printf("Number of Answers: %d\n", dns.number_of_answers_rrs)
	fmt.Printf("Number of Authorities: %d\n", dns.number_of_authority_rrs)
	fmt.Printf("Number of Additional: %d\n", dns.number_of_additional_rrs)
	fmt.Println("BODY INFO")
	for _, question := range dns.questions {
		fmt.Printf("Hostname to query: %s\n", question.name)
		fmt.Printf("Query Type: %d\n", question.qtype)
		fmt.Printf("Query Class: %d\n", question.class)
	}

	for _, answers := range dns.resource_records {
		fmt.Printf("Response Name: %s\n", answers.name)
		fmt.Printf("RR Type: %d\n", answers.rrtype)
		fmt.Printf("Response Class: %d\n", answers.class)
		fmt.Printf("Response TTL: %d\n", answers.ttl)
		fmt.Printf("Response Length: %d\n", answers.rdlength)
		if answers.rrtype == 1 {
			fmt.Printf("%d.%d.%d.%d\n", answers.rdata[0], answers.rdata[1], answers.rdata[2], answers.rdata[3])
		} else {
			fmt.Printf("Response Data: %s\n", string(answers.rdata))
		}
	}
}
