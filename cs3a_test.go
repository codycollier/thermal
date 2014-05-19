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

func Test3aGenerateLineEncryptionKey(t *testing.T) {

	var lineSecret [32]byte
	var lineSecret2 [32]byte
	var localLineId [16]byte
	var remoteLineId [16]byte

	rand.Reader.Read(lineSecret[:])
	rand.Reader.Read(localLineId[:])
	rand.Reader.Read(remoteLineId[:])

	cset := new(cs3a)
	cset.initialize()

	testkey1 := cset.generateLineEncryptionKey(&lineSecret, &localLineId, &remoteLineId)
	testkey2 := cset.generateLineEncryptionKey(&lineSecret, &localLineId, &remoteLineId)
	if testkey1 != testkey2 {
		t.Log("Repeat call to generate should produce same key")
		t.Fail()
	}

	testkey3 := cset.generateLineEncryptionKey(&lineSecret2, &remoteLineId, &localLineId)
	if testkey3 == testkey1 {
		t.Log("New line param order should generate should produce new key")
		t.Fail()
	}

	rand.Reader.Read(lineSecret2[:])
	testkey4 := cset.generateLineEncryptionKey(&lineSecret2, &localLineId, &remoteLineId)
	if testkey4 == testkey1 {
		t.Log("New line secret should produce new key")
		t.Fail()
	}

}

func Test3aGenerateLineDecryptionKey(t *testing.T) {

	var lineSecret [32]byte
	var lineSecret2 [32]byte
	var localLineId [16]byte
	var remoteLineId [16]byte

	rand.Reader.Read(lineSecret[:])
	rand.Reader.Read(localLineId[:])
	rand.Reader.Read(remoteLineId[:])

	cset := new(cs3a)
	cset.initialize()

	testkey1 := cset.generateLineDecryptionKey(&lineSecret, &localLineId, &remoteLineId)
	testkey2 := cset.generateLineDecryptionKey(&lineSecret, &localLineId, &remoteLineId)
	if testkey1 != testkey2 {
		t.Log("Repeat call to generate should produce same key")
		t.Fail()
	}

	testkey3 := cset.generateLineDecryptionKey(&lineSecret2, &remoteLineId, &localLineId)
	if testkey3 == testkey1 {
		t.Log("New line param order should generate should produce new key")
		t.Fail()
	}

	rand.Reader.Read(lineSecret2[:])
	testkey4 := cset.generateLineDecryptionKey(&lineSecret2, &localLineId, &remoteLineId)
	if testkey4 == testkey1 {
		t.Log("New line secret should produce new key")
		t.Fail()
	}

}

func Test3aEncryptLinePacket(t *testing.T) {

	var localLineSecret [32]byte
	var localLineId [16]byte
	var remoteLineId [16]byte

	rand.Reader.Read(localLineSecret[:])
	rand.Reader.Read(localLineId[:])
	rand.Reader.Read(remoteLineId[:])

	cset := new(cs3a)
	cset.initialize()

	packet := []byte("This is a channel packet inside the line packet")
	lineEncryptionKey := cset.generateLineEncryptionKey(&localLineSecret, &localLineId, &remoteLineId)
	linePacketBody, err := cset.encryptLinePacket(packet, &lineEncryptionKey)
	if err != nil {
		t.Logf("Error encrypting line packet")
		t.Fail()
	}

	if len(linePacketBody) < 25 {
		t.Logf("linePacketBody does not have enough bytes")
		t.Fail()
	}

}

func Test3aEncryptAndDecryptLinePacket(t *testing.T) {

	var localLineSecret [32]byte
	var localLineId [16]byte
	var remoteLineId [16]byte

	rand.Reader.Read(localLineSecret[:])
	rand.Reader.Read(localLineId[:])
	rand.Reader.Read(remoteLineId[:])

	cset := new(cs3a)
	cset.initialize()

	lineEncryptionKey := cset.generateLineEncryptionKey(&localLineSecret, &localLineId, &remoteLineId)
	lineDecryptionKey := cset.generateLineDecryptionKey(&localLineSecret, &remoteLineId, &localLineId)

	originalPacket := []byte("This is a channel packet inside the line packet")
	linePacketBody, err := cset.encryptLinePacket(originalPacket, &lineEncryptionKey)
	if err != nil {
		t.Logf("Error encrypting line packet")
		t.Fail()
	}
	returnedPacket, err := cset.decryptLinePacket(linePacketBody, &lineDecryptionKey)
	if err != nil {
		t.Logf("Error decrypting line packet")
		t.Fail()
	}

	t.Logf("originalPacket: %x\n", originalPacket)
	t.Logf("returnedPacket: %x\n", returnedPacket)
	t.Logf("\n")

	if !bytes.Equal(originalPacket, returnedPacket) {
		t.Logf("The packet was corrupted during encryption and decryption")
		t.Fail()
	}

}
