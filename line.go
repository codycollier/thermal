package thermal

type lineHalf struct {
	id           [16]byte
	at           int
	sharedSecret [32]byte
}

type lineSession struct {
	local         lineHalf
	remote        lineHalf
	encryptionKey [32]byte
	decryptionKey [32]byte
}

type lineMap map[string]lineSession
