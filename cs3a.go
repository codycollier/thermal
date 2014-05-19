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

/*
--------------------------------------------------------------------------------
The CS3a Open Packet Handshake
---------------------------------

The cs3a encryption & decryption of an inner open packet, from the perspective
of the sender and receiver:

Sender:
box.Precompute(senderLineSecret, receiverPublicKey, senderLinePrivateKey)
secretbox.Seal(encInnerPacket, packet, &nonce, senderLineSecret)

Receiver:
box.Precompute(&senderLineSecret, &senderLinePublicKey, &receiverPrivateKey)
secretbox.Open(packet, encInnerPacket, &nonce, &senderLineSecret)

The sender/receiver context helps to highlight the public/private key pairings.

However, the following functions use the context of local switch instance
and remote switch instance instead.
--------------------------------------------------------------------------------
*/

// encryptOpenPacket returns an assembled open packet body and a local line shared secret
func (cs *cs3a) encryptOpenPacket(packet []byte, remotePublicKey *[32]byte) (openPacketBody []byte, localLineSecret [32]byte, err error) {

	// switch key pair
	// cs.publicKey and cs.privateKey should already be populated

	// temporary line key pair
	linePublicKey, linePrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Println("Error generating NaCl keypair for line")
		return openPacketBody, localLineSecret, err
	}

	// Encrypt the inner packet
	var nonce [24]byte
	var encInnerPacket []byte

	box.Precompute(&localLineSecret, remotePublicKey, linePrivateKey)
	encInnerPacket = secretbox.Seal(encInnerPacket, packet, &nonce, &localLineSecret)

	// Generate the mac and assemble the body for the outer packet
	// <mac><local-line-public-key><encrypted-inner-packet-data>
	var macKey [32]byte
	var mac [16]byte
	var openPacketData []byte

	box.Precompute(&macKey, remotePublicKey, &cs.privateKey)
	openPacketData = append(openPacketData, linePublicKey[:]...)
	openPacketData = append(openPacketData, encInnerPacket...)
	poly1305.Sum(&mac, openPacketData, &macKey)

	openPacketBody = append(openPacketBody, mac[:]...)
	openPacketBody = append(openPacketBody, openPacketData...)

	return openPacketBody, localLineSecret, nil

}

// decryptOpenPacket returns an unencrypted inner open packet and a remote line shared secret
func (cs *cs3a) decryptOpenPacket(openPacketBody []byte, remotePublicKey *[32]byte) (packet []byte, remoteLineSecret [32]byte, err error) {

	// switch key pair
	// cs.publicKey and cs.privateKey should already be populated

	// Unpack the outer packet body
	// <mac><remote-line-public-key><encrypted-inner-packet-data>
	var mac [16]byte
	var remoteLinePublicKey [32]byte
	var encInnerPacket []byte
	var openPacketData []byte

	copy(mac[:], openPacketBody[:16])
	copy(remoteLinePublicKey[:], openPacketBody[16:48])
	encInnerPacket = append(encInnerPacket, openPacketBody[48:]...)
	openPacketData = append(openPacketData, remoteLinePublicKey[:]...)
	openPacketData = append(openPacketData, encInnerPacket...)

	// Verify the mac
	var authenticated bool
	var macKey [32]byte

	box.Precompute(&macKey, remotePublicKey, &cs.privateKey)
	authenticated = poly1305.Verify(&mac, openPacketData, &macKey)
	if !authenticated {
		msg := "Incoming open packet failed MAC authentication"
		log.Println(msg)
		err = fmt.Errorf(msg)
		return packet, remoteLineSecret, err
	}

	// Decrypt the inner packet
	var nonce [24]byte

	box.Precompute(&remoteLineSecret, &remoteLinePublicKey, &cs.privateKey)
	packet, success := secretbox.Open(packet, encInnerPacket, &nonce, &remoteLineSecret)
	if !success {
		err := fmt.Errorf("Error opening the secretbox")
		return packet, remoteLineSecret, err
	}

	log.Printf("packet: %x\n", packet)
	return packet, remoteLineSecret, nil
}

// generateLineEncryptionKey returns a key suitable for outgoing line packet encryption,
// in the form of sha256(local-line-secret, local-line-id, local-line-id)
func (cs *cs3a) generateLineEncryptionKey(localLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) (key [32]byte) {

	hash256 := sha256.New()
	hash256.Write(localLineSecret[:])
	hash256.Write(localLineId[:])
	hash256.Write(remoteLineId[:])
	keyHash := hash256.Sum(nil)

	copy(key[:], keyHash[:])

	return key
}

// generateLineDecryptionKey returns a key suitable for incoming line packet decryption,
// in the form of sha256(remote-line-secret, remote-line-id, local-line-id)
func (cs *cs3a) generateLineDecryptionKey(remoteLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) (key [32]byte) {

	hash256 := sha256.New()
	hash256.Write(remoteLineSecret[:])
	hash256.Write(remoteLineId[:])
	hash256.Write(localLineId[:])
	keyBin := hash256.Sum(nil)

	copy(key[:], keyBin[:])

	return key
}

// encryptLinePacket returns an encrypted and assembled outer line packet body
// in the form <nonce><secretbox-ciphertext>
func (cs *cs3a) encryptLinePacket(packet []byte, lineEncryptionKey *[32]byte) (linePacketBody []byte) {
	var nonce [24]byte
	var cipherText = make([]byte, 0)

	rand.Reader.Read(nonce[:])
	secretbox.Seal(cipherText, packet, &nonce, lineEncryptionKey)
	linePacketBody = append(nonce[:], cipherText...)

	return linePacketBody
}

// decryptLinePacket returns a decrypted inner line packet
// by disassembling and decrypting an outer packet body in the form <nonce><secretbox-ciphertext>
func (cs *cs3a) decryptLinePacket(linePacketBody []byte, lineDecryptionKey *[32]byte) (packet []byte) {
	var nonce [24]byte
	copy(nonce[:], linePacketBody[:24])
	linePacketData := linePacketBody[24:]
	secretbox.Open(packet, linePacketData, &nonce, lineDecryptionKey)
	return packet
}
