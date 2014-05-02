package thermal

import (
//"code.google.com/p/go.crypto/nacl/box"
//"code.google.com/p/go.crypto/nacl/secretbox"
//"crypto/rand"
//"crypto/sha256"
)

type cs3a cipherSetPlugin

func (cs *cs3a) init() {
	cs.id = "cs3a"

}
