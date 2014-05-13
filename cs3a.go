package thermal

import (
	//"code.google.com/p/go.crypto/nacl/secretbox"
	"code.google.com/p/go.crypto/nacl/box"
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

func (cs *cs3a) encryptOpenPacket(packet []byte) (csdata []byte) {

	//------------------------------------------------------
	// Setup the keys
	//------------------------------------------------------
	// . existing switch key pair
	// . new line key pair?	// return?  store in local map?

	//------------------------------------------------------
	// Encrypt the inner packet
	//------------------------------------------------------
	// . 24 byte nonce of 0
	// nonce = make([]byte, 24)

	// . generate the shared secret
	// crypto_box(agreedKey, receiver.pubkey, sender.line.prvkey)

	// . use secretbox to encrypt the inner packet
	// encInnerPacket = crypto_secretbox(data, nonce, agreedKey)

	//------------------------------------------------------
	// Generate the hmac and assemble the outer packet
	//------------------------------------------------------
	// . <open-HMAC><sender-line-public-key><encrypted-inner-packet-data>
	// .		csdata.00-32  == auth - onetimeauth
	// .		csdata.33-64  == line_key - senders line level public key
	// .		csdata.rest   == inner_ciphertext (encrypted open packet)

	// . assemble part of the outer packet
	// outerPacketData = sender_linekey.public + encInnerPacket

	// . generate the macKey for use in generating the hmac
	// crypto_box(macKey, receiver.pubkey, sender.prvkey)

	// . generate the hmac
	// hmac = crypto_onetimeauth(encInnerPacket, macKey)
	// var hmac
	// crypt.poly1305.sum(hmac, encInnerPacket, macKey)

	// . assemble the rest of the outer packet BODY
	// csdata = hmac + outerPacketData
	//
	return csdata

}

func (cs *cs3a) decryptOpenPacket(csdata []byte) (packet []byte) {
	return packet
}

func (cs *cs3a) encryptLinePacket(packet []byte) (csdata []byte) {
	return csdata
}

func (cs *cs3a) decryptLinePacket(csdata []byte) (packet []byte) {
	return packet
}
