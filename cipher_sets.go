package thermal

// The cipherSet interface defines behavior expected from a cipher set plugin
type cipherSet interface {
	initialize() error
	fingerprint() (string, string)

	encryptOpenPacket(packet []byte, receiverPublicKey *[32]byte) (openPacketBody []byte, lineSharedSecret *[32]byte, err error)
	decryptOpenPacket(openPacketBody []byte, senderPublicKey *[32]byte) (packet []byte, lineSharedSecret [32]byte, err error)

	generateLineEncryptionKey(lineSharedSecret *[32]byte, localLineId, remoteLineId *[16]byte) [32]byte
	generateLineDecryptionKey(lineSharedSecret *[32]byte, localLineId, remoteLineId *[16]byte) [32]byte

	encryptLinePacket(packet []byte, lineEncryptionKey *[32]byte) (linePacketBody []byte)
	decryptLinePacket(linePacketBody []byte, lineDecryptionKey *[32]byte) (packet []byte)
}

// A cipherPack holds a group of cipher sets
type cipherPack map[string]cipherSet
