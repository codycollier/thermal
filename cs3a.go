package thermal

import (
	"code.google.com/p/go.crypto/nacl/box"
	"code.google.com/p/go.crypto/nacl/secretbox"
	"code.google.com/p/go.crypto/poly1305"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
)

// cs3a is an implementation of the NaCl based cipher set 3a
type cs3a struct {
	id             string
	fingerprintBin []byte
	fingerprintHex string
	publicKey      [32]byte
	privateKey     [32]byte
}

// initialize generates a key pair and sets up the cipher set
func (cs *cs3a) initialize() error {

	// generate the key pair
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Println("Error generating NaCl keypair in cs3a initialization")
		return err
	}

	// generate the fingerprint hash
	hash256 := sha256.New()
	hash256.Write(publicKey[:])
	fingerprintBin := hash256.Sum(nil)
	fingerprintHex := fmt.Sprintf("%x", fingerprintBin)

	// initialize the struct
	cs.id = "cs3a"
	cs.fingerprintBin = fingerprintBin
	cs.fingerprintHex = fingerprintHex
	cs.publicKey = *publicKey
	cs.privateKey = *privateKey

	return nil

}

func (cs *cs3a) String() string {
	return fmt.Sprintf("%s: %x", cs.id, cs.fingerprint)
}

func (cs *cs3a) csid() string {
	return cs.id
}

// fingerprint returns the csid and fingerprint for use in a 'parts' set
func (cs *cs3a) fingerprint() (string, string) {
	return cs.id, cs.fingerprintHex
}

func (cs *cs3a) encryptOpenPacket(packet []byte, receiverPublicKey *[32]byte) (openPacketBody []byte, err error) {

	// todo - return/store the egress line shared key for use in line packet encryption

	// switch key pair
	// cs.publicKey and cs.privateKey should already be populated

	// line key pair
	linePublicKey, linePrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Println("Error generating NaCl keypair for line")
		return openPacketBody, err
	}

	// Encrypt the inner packet
	var nonce [24]byte
	var lineSharedKey [32]byte
	var encInnerPacket []byte

	box.Precompute(&lineSharedKey, receiverPublicKey, linePrivateKey)
	secretbox.Seal(encInnerPacket, packet, &nonce, &lineSharedKey)

	// Generate the mac and assemble the body for the outer packet
	// <mac><sender-line-public-key><encrypted-inner-packet-data>
	var macKey [32]byte
	var mac [16]byte
	var openPacketData []byte

	box.Precompute(&macKey, receiverPublicKey, &cs.privateKey)
	openPacketData = append(linePublicKey[:], encInnerPacket...)
	poly1305.Sum(&mac, openPacketData, &macKey)
	openPacketBody = append(mac[:], openPacketData...)

	return openPacketBody, nil

}

func (cs *cs3a) decryptOpenPacket(openPacketBody []byte, senderPublicKey *[32]byte) (packet []byte, err error) {

	// todo - return/store the ingress line shared key for use in line packet decryption

	// switch key pair
	// cs.publicKey and cs.privateKey should already be populated

	// Unpack the outer packet body
	// <mac><sender-line-public-key><encrypted-inner-packet-data>
	var mac [16]byte
	var senderLinePublicKey [32]byte
	var encInnerPacket []byte
	var openPacketData []byte

	copy(mac[:], openPacketBody[:16])
	copy(senderLinePublicKey[:], openPacketBody[16:48])
	copy(encInnerPacket[:], openPacketBody[48:])
	openPacketData = append(senderLinePublicKey[:], encInnerPacket...)

	// Verify the mac
	var authenticated bool
	var macKey [32]byte

	box.Precompute(&macKey, senderPublicKey, &cs.privateKey)
	authenticated = poly1305.Verify(&mac, openPacketData, &macKey)
	if !authenticated {
		msg := "Incoming open packet failed MAC authentication"
		log.Println(msg)
		err = fmt.Errorf(msg)
		return packet, err
	}

	// Decrypt the inner packet
	var nonce [24]byte
	var lineSharedKey [32]byte

	box.Precompute(&lineSharedKey, &senderLinePublicKey, &cs.privateKey)
	secretbox.Open(packet, encInnerPacket, &nonce, &lineSharedKey)

	return packet, nil
}

func (cs *cs3a) encryptLinePacket(packet []byte) (linePacketBody []byte) {
	return linePacketBody
}

func (cs *cs3a) decryptLinePacket(linePacketBody []byte) (packet []byte) {
	return packet
}
