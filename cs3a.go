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
	fingerprint []byte
	private_key [32]byte
	public_key  [32]byte
}

// init will initialize an empty cs3a
func (cs *cs3a) init() {

	pubkey, prvkey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Panicf("Error: %s\n", err)
	}

	hash256 := sha256.New()
	hash256.Write(pubkey[:])
	fingerprint := hash256.Sum(nil)

	cs.id = "cs3a"
	cs.fingerprint = fingerprint
	cs.public_key = *pubkey
	cs.private_key = *prvkey

}

func (cs *cs3a) String() string {
	return fmt.Sprintf("%s: %x", cs.id, cs.fingerprint)
}

// csid will return the id of the cipher set
func (cs *cs3a) csid() string {
	return cs.id
}

// return the 'parts' for the cipherset
func (cs *cs3a) parts() (string, []byte) {
	return cs.id, cs.fingerprint
}
