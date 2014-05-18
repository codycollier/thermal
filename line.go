package thermal

type line struct {
	id           [16]byte
	at           int
	sharedSecret [32]byte
}

type lineSession struct {
	localLine     line
	remoteLine    line
	encryptionKey [32]byte
	decryptionKey [32]byte
}

type lineMap map[string]lineSession
