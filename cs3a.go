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

// init will generate a key pair and initialize the cipher set
func (cs *cs3a) init() error {

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

// csid will return the id of the cipher set
func (cs *cs3a) csid() string {
	return cs.id
}

// fingerprint will return the csid and fingerprint for use in a 'parts' set
func (cs *cs3a) fingerprint() (string, string) {
	return cs.id, cs.fingerprintHex
}
