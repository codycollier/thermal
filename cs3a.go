package thermal

import (
	//"code.google.com/p/go.crypto/nacl/secretbox"
	"code.google.com/p/go.crypto/nacl/box"
	"crypto/rand"
	"crypto/sha256"
	"log"
)

type cs3a cipherSetPlugin

func (cs *cs3a) init() {

	pubkey, prvkey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	hash256 := sha256.New()
	hash256.Write(pubkey[:])
	fingerprint := hash256.Sum(nil)

	cs.id = "cs3a"
	cs.fingerprint = fingerprint
	cs.public_key = *pubkey
	cs.private_key = *prvkey
}
