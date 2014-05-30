package thermal

import (
	"fmt"
)

// The cipherSet interface defines behavior expected from a cipher set plugin
type cipherSet interface {
	initialize() error
	csid() [1]byte
	fingerprint() (string, string)
	pubKey() *[32]byte

	encryptOpenPacketBody(packet []byte, remotePublicKey *[32]byte) (openPacketBody []byte, localLineSecret [32]byte, err error)
	decryptOpenPacketBody(openPacketBody []byte, remotePublicKey *[32]byte) (packet []byte, remoteLineSecret [32]byte, err error)

	generateLineEncryptionKey(localLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) [32]byte
	generateLineDecryptionKey(remoteLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) [32]byte

	encryptLinePacketBody(packet []byte, lineEncryptionKey *[32]byte) (linePacketBody []byte, err error)
	decryptLinePacketBody(linePacketBody []byte, lineDecryptionKey *[32]byte) (packet []byte, err error)
}

// A cipherPack holds a group of cipher sets
type cipherPack map[string]cipherSet

// extractParts will pull out the hashname parts from a cipherPack
func extractParts(cpack *cipherPack) (map[string]string, error) {

	if len(*cpack) == 0 {
		return nil, fmt.Errorf("cipherPack is empty")
	}

	parts := make(map[string]string)

	for _, cset := range *cpack {
		csid, fingerprint := cset.fingerprint()

		if fingerprint == "" {
			return nil, fmt.Errorf("fingerprint for %s is empty", csid)
		}

		parts[csid] = fingerprint
	}

	return parts, nil
}
