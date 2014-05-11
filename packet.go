package thermal

// decodedPacket holds the pieces of a decoded telehash packet
// packet encoding:
//  <length-of-head>[HEAD][BODY]
type decodedPacket struct {
	headLength int
	head       []byte
	json       string
	bodyLength int
	body       []byte
}

// decodePacket accepts a telehash packet and breaks it in to components
func decodePacket(packet []byte) (decodedPacket, error) {
}

// encodePacket accepts head and body payloads and encodes them to a packet
func encodePacket(json string, body []byte) (packet []byte, err error) {
}
