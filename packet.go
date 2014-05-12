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
	var head []byte = nil
	var bodyLength int
	var body []byte = nil
	var json string = ""
	var err error = nil

	if len(packetBytes) < 3 {
		packet = decodedPacket{}
		err = fmt.Errorf("Packet must be at least three bytes long\n")
		return packet, err
	}

	// The head-length will be the first two bytes of the packet (network byte order / big endian)
	headLength = int(binary.BigEndian.Uint16(packetBytes[:2]))

	bodyLength = len(packetBytes) - headLength - 2

	if headLength > 0 {
		head = packetBytes[2 : 2+headLength]
		json = string(head)
	}

	if bodyLength > 0 {
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

	var packet = make([]byte, 0)
	var err error

	var headLength = make([]byte, 2)
	var head []byte = nil

	if json != "" {
		head = []byte(json)
	}

	// The head-length will be the first two bytes of the packet (network byte order / big endian)
	binary.BigEndian.PutUint16(headLength, uint16(len(head)))

	// Assemble the packet
	packet = append(packet, headLength...)
	packet = append(packet, head...)
	packet = append(packet, body...)

	return packet, err

}
