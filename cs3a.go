package thermal

import (
	"code.google.com/p/go.crypto/nacl/box"
	"code.google.com/p/go.crypto/nacl/secretbox"
	"code.google.com/p/go.crypto/poly1305"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

// cs3a is an implementation of the NaCl based cipher set 3a
type cs3a struct {
	id             [1]byte
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

	// generate the single byte representation of the cipher set
	csid_byte, _ := hex.DecodeString("3a")

	// initialize the struct
	copy(cs.id[:], csid_byte[:])
	cs.fingerprintBin = fingerprintBin
	cs.publicKey = *publicKey
	cs.privateKey = *privateKey

	return nil

}

func (cs *cs3a) String() string {
	return fmt.Sprintf("%x: %x", cs.id, cs.fingerprintBin)
}

func (cs *cs3a) csid() [1]byte {
	return cs.id
}

// fingerprint returns the csid and fingerprint for use in a 'parts' set
func (cs *cs3a) fingerprint() (string, string) {
	return fmt.Sprintf("%x", cs.id), fmt.Sprintf("%x", cs.fingerprintBin)
}

// pubKey returns the cipher set public key
func (cs *cs3a) pubKey() *[32]byte {
	return &cs.publicKey
}

/*
--------------------------------------------------------------------------------
gob encoding/decoding





--------------------------------------------------------------------------------
*/

// GobEncode implements the GobEncoder interface and allows for persisting the cipherset
func (cs *cs3a) GobEncode() ([]byte, error) {

	var encoded_cset []byte

	encoded_cset = append(encoded_cset, cs.id[:]...)
	encoded_cset = append(encoded_cset, cs.publicKey[:]...)
	encoded_cset = append(encoded_cset, cs.privateKey[:]...)

	return encoded_cset, nil

}

/*
--------------------------------------------------------------------------------
The CS3a Open Packet Handshake

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
func (cs *cs3a) encryptOpenPacketBody(packet []byte, remotePublicKey *[32]byte) (openPacketBody []byte, localLineSecret [32]byte, err error) {

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
func (cs *cs3a) decryptOpenPacketBody(openPacketBody []byte, remotePublicKey *[32]byte) (packet []byte, remoteLineSecret [32]byte, err error) {

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

	//log.Printf("packet: %x\n", packet)
	return packet, remoteLineSecret, nil
}

/*
--------------------------------------------------------------------------------
The CS3a Line Packet encryption

--------------------------------------------------------------------------------
*/

// generateLineEncryptionKey returns a key suitable for outgoing line packet encryption
func (cs *cs3a) generateLineEncryptionKey(localLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) (key [32]byte) {

	// sha256(local-line-secret, local-line-id, local-line-id)
	hash256 := sha256.New()
	hash256.Write(localLineSecret[:])
	hash256.Write(localLineId[:])
	hash256.Write(remoteLineId[:])
	keyHash := hash256.Sum(nil)

	copy(key[:], keyHash[:])

	return key
}

// generateLineDecryptionKey returns a key suitable for incoming line packet decryption
func (cs *cs3a) generateLineDecryptionKey(remoteLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) (key [32]byte) {

	// sha256(remote-line-secret, remote-line-id, local-line-id)
	hash256 := sha256.New()
	hash256.Write(remoteLineSecret[:])
	hash256.Write(remoteLineId[:])
	hash256.Write(localLineId[:])
	keyBin := hash256.Sum(nil)

	copy(key[:], keyBin[:])

	return key
}

// encryptLinePacket encrypts a channel packet and builds a line packet body
func (cs *cs3a) encryptLinePacketBody(packet []byte, lineEncryptionKey *[32]byte) (linePacketBody []byte, err error) {
	var nonce [24]byte
	var linePacketData []byte

	// encrypt the inner channel packet
	rand.Reader.Read(nonce[:])
	linePacketData = secretbox.Seal(linePacketData, packet, &nonce, lineEncryptionKey)

	// assemble: <nonce><secretbox-ciphertext>
	linePacketBody = append(nonce[:], linePacketData[:]...)

	return linePacketBody, nil
}

// decryptLinePacket returns a decrypted channel packet from a line packet
func (cs *cs3a) decryptLinePacketBody(linePacketBody []byte, lineDecryptionKey *[32]byte) (packet []byte, err error) {
	var nonce [24]byte

	// disassemble: <nonce><secretbox-ciphertext>
	copy(nonce[:], linePacketBody[:24])
	linePacketData := linePacketBody[24:]

	// decrypt the inner channel packet
	packet, success := secretbox.Open(packet, linePacketData, &nonce, lineDecryptionKey)
	if !success {
		err := fmt.Errorf("Error opening the secretbox")
		return packet, err
	}

	return packet, nil
}
