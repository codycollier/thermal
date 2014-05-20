package thermal

import (
	"crypto/rand"
	"fmt"
)

type lineHalf struct {
	id     [16]byte
	at     int
	secret [32]byte
}

type lineSession struct {
	local         lineHalf
	remote        lineHalf
	encryptionKey [32]byte
	decryptionKey [32]byte
}

type lineMap map[string]lineSession

// generateLineId returns a random 16 char hex encoded string
func generateLineId() string {
	var idBin [8]byte
	rand.Reader.Read(idBin[:])
	idHex := fmt.Sprintf("%x", idBin)
	return idHex
}
