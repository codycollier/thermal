package thermal

// The cipherSet interface defines behavior expected from a cipher set plugin
type cipherSet interface {
	initialize() error
	fingerprint() (string, string)

	encryptOpenPacket(packet []byte, remotePublicKey *[32]byte) (openPacketBody []byte, localLineSecret [32]byte, err error)
	decryptOpenPacket(openPacketBody []byte, remotePublicKey *[32]byte) (packet []byte, remoteLineSecret [32]byte, err error)

	generateLineEncryptionKey(localLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) [32]byte
	generateLineDecryptionKey(remoteLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) [32]byte

	encryptLinePacket(packet []byte, lineEncryptionKey *[32]byte) (linePacketBody []byte)
	decryptLinePacket(linePacketBody []byte, lineDecryptionKey *[32]byte) (packet []byte)
}

// A cipherPack holds a group of cipher sets
type cipherPack map[string]cipherSet
