package thermal

import (
	"encoding/binary"
	"fmt"
)

// A decodedPacket holds the pieces of a decoded telehash packet
// packet encoding:
//  <length-of-head>[HEAD][BODY]
type decodedPacket struct {
	headLength int
	head       []byte
	bodyLength int
	body       []byte
	json       string
}

// decodePacket accepts a telehash packet and breaks it in to components
func decodePacket(packetBytes []byte) (decodedPacket, error) {

	var packet decodedPacket
	var headLength int
	var head []byte
	var bodyLength int
	var body []byte
	var json string
	var err error = nil

	if len(packetBytes) < 3 {
		packet = decodedPacket{}
		err = fmt.Errorf("Packet must be at least three bytes long\n")
		return packet, err
	}

	headLength = int(binary.BigEndian.Uint16(packetBytes[:2]))

	bodyLength = len(packetBytes) - headLength - 2

	if headLength == 0 {
		head = nil
		json = ""
	} else {
		head = packetBytes[2 : 2+headLength]
		json = string(head)
	}

	if bodyLength == 0 {
		body = nil
	} else {
		body = packetBytes[2+headLength:]
	}

	packet = decodedPacket{
		headLength: headLength,
		bodyLength: bodyLength,
		head:       head,
		body:       body,
		json:       json,
	}

	return packet, nil
}

// encodePacket accepts json and body payloads and encodes them to a packet
func encodePacket(json string, body []byte) ([]byte, error) {

	var headLength = make([]byte, 2)
	var packet, head []byte
	var err error
	packet = make([]byte, 0)

	if json != "" {
		head = []byte(json)
	} else {
		head = nil
	}

	// The head-length will be the first two bytes of the packet (network byte order / big endian)
	/*
		headBytes := new(bytes.Buffer)
		err = binary.Write(headBytes, binary.BigEndian, len(head))
		if err != nil {
			return nil, fmt.Errorf("Error writing head length to packet (err: %s", err)
		}
		headLength = []byte(headBytes.Bytes())
	*/
	binary.BigEndian.PutUint16(headLength, uint16(len(head)))

	// Assemble the packet
	packet = append(packet, headLength...)
	packet = append(packet, head...)
	packet = append(packet, body...)

	return packet, err

}
