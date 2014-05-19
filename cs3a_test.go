package thermal

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func Test3aInit(t *testing.T) {

	cset := new(cs3a)
	err := cset.initialize()
	if err != nil {
		t.Log("Error in cs3a initialization")
		t.Logf("err: %s\n", err)
		t.Fail()
	}

}

func Test3aCsid(t *testing.T) {
	cset := new(cs3a)
	cset.initialize()
	csid := cset.csid()
	if csid != "cs3a" {
		t.Fail()
	}
}

func Test3aFingerprint(t *testing.T) {
	cset := new(cs3a)
	cset.initialize()
	csid, fingerprint := cset.fingerprint()
	if csid != "cs3a" {
		t.Fail()
	}
	if len(fingerprint) != 64 {
		t.Fail()
	}

}

func Test3aEncryptOpenPacket(t *testing.T) {
	var remotePublicKey [32]byte
	rand.Reader.Read(remotePublicKey[:])

	originalPacket := []byte("This is an internal open packet")
	cset := new(cs3a)
	cset.initialize()

	openPacketBody, lineSecretA, err := cset.encryptOpenPacket(originalPacket, &remotePublicKey)
	if err != nil {
		t.Logf("Error: %s", err)
		t.Fail()
	}

	if len(lineSecretA) != 32 {
		t.Logf("line secret is not expected length: %x", lineSecretA)
		t.Fail()
	}

	if len(openPacketBody) < 49 {
		t.Logf("openPacketBody does not have enough bytes")
		t.Fail()
	}

}

func Test3aEncryptAndDecryptOpenPacket(t *testing.T) {

	var err error
	originalPacket := []byte("This is an internal open packet")

	csetSender := new(cs3a)
	err = csetSender.initialize()

	csetReceiver := new(cs3a)
	err = csetReceiver.initialize()

	openPacketBody, lineSecretA, err := csetSender.encryptOpenPacket(originalPacket, &csetReceiver.publicKey)
	if err != nil {
		t.Logf("Error: %s", err)
		t.Fail()
	}
	returnedPacket, lineSecretB, err := csetReceiver.decryptOpenPacket(openPacketBody, &csetSender.publicKey)
	if err != nil {
		t.Logf("Error: %s", err)
		t.Fail()
	}

	t.Logf("originalPacket: %s", originalPacket)
	t.Logf("returnedPacket: %s", returnedPacket)
	if !bytes.Equal(returnedPacket, originalPacket) {
		t.Log("The returned packet did not match the original packet")
		t.Fail()
	}

	t.Logf("line secret A: %x", lineSecretA)
	t.Logf("line secret B: %x", lineSecretB)
	if lineSecretA != lineSecretB {
		t.Log("The line shared secret does not match")
		t.Fail()
	}

}

/*
func (cs *cs3a) encryptOpenPacket(packet []byte, remotePublicKey *[32]byte) (openPacketBody []byte, localLineSecret [32]byte, err error) {
func (cs *cs3a) decryptOpenPacket(openPacketBody []byte, remotePublicKey *[32]byte) (packet []byte, remoteLineSecret [32]byte, err error) {
func (cs *cs3a) generateLineEncryptionKey(localLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) (key [32]byte) {
func (cs *cs3a) generateLineDecryptionKey(remoteLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) (key [32]byte) {
func (cs *cs3a) encryptLinePacket(packet []byte, lineEncryptionKey *[32]byte) (linePacketBody []byte) {
func (cs *cs3a) decryptLinePacket(linePacketBody []byte, lineDecryptionKey *[32]byte) (packet []byte) {
*/
