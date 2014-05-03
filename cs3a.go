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
	id          string
	fingerprint string
	publicKey   [32]byte
	privateKey  [32]byte
}

// init will generate a key pair and initialize the cipher set
func (cs *cs3a) init() {

	// generate the key pair
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Println("Error generating NaCl keypair in cs3a initialization")
		log.Panicf("Error: %s\n", err)
	}

	// generate the fingerprint hash
	hash256 := sha256.New()
	hash256.Write(publicKey[:])
	fingerprintBytes := hash256.Sum(nil)
	fingerprintHex := fmt.Sprintf("%x", fingerprintBytes)

	// initialize the struct
	cs.id = "cs3a"
	cs.fingerprint = fingerprintHex
	cs.publicKey = *publicKey
	cs.privateKey = *privateKey

}

func (cs *cs3a) String() string {
	return fmt.Sprintf("%s: %x", cs.id, cs.fingerprint)
}

// csid will return the id of the cipher set
func (cs *cs3a) csid() string {
	return cs.id
}

// parts will return the telehash defined 'parts' for the cipherset
func (cs *cs3a) parts() (string, string) {
	return cs.id, cs.fingerprint
}
