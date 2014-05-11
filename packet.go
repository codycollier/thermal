package thermal

import (
	"encoding/json"
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

// Should the json decoding be handled here as well?  Or in a channel/type
// specific area?  If here, create a map of json/channel types and their
// structs?  It might be good to have all the json schema here in the same
// place.

// decodePacket accepts a telehash packet and breaks it in to components
func decodePacket(packet []byte) (decodedPacket, error) {
}

// encodePacket accepts json and body payloads and encodes them to a packet
func encodePacket(json string, body []byte) ([]byte, error) {
}
