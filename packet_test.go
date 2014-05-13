package thermal

import (
	"bytes"
	"testing"
)

func TestBasicPacketEncodeDecode(t *testing.T) {

	expected_json := "{\"color\": \"green\"}"
	expected_body := []byte("This is a test body")

	encPacket, err := encodePacket(expected_json, expected_body)
	if err != nil {
		t.Log("Error encoding packet (err: %s)\n", err)
		t.Fail()
	}

	packet, err := decodePacket(encPacket)
	if err != nil {
		t.Log("Error decoding packet (err: %s)\n", err)
		t.Fail()
	}

	t.Logf("packet.json: %s\n", packet.json)
	t.Logf("expected_json: %s\n", expected_json)
	t.Logf("packet.body: %s\n", packet.body)
	t.Logf("expected_body: %s\n", expected_body)

	if packet.json != expected_json {
		t.Log("json was found to be corrupted by encoding / decoding process")
		t.Fail()
	}

	if bytes.Compare(packet.body, expected_body) != 0 {
		t.Log("body was found to be corrupted by encoding / decoding process")
		t.Fail()
	}

}