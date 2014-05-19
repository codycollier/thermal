package thermal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"time"
)

// cs2a is an implementation of the cipher set 2a
type cs2a struct {
	id             string
	fingerprintBin []byte
	fingerprintHex string
	publicKey      []byte
	privateKey     rsa.PrivateKey
	certificate    x509.Certificate
}

func (cs *cs2a) String() string {
	return fmt.Sprintf("%s: %x", cs.id, cs.fingerprint)
}

// csid will return the id of the cipher set
func (cs *cs2a) csid() string {
	return cs.id
}

// init will generate a key pair and initialize the cipher set
func (cs *cs2a) initialize() error {

	// generate the rsa-2048 key pair
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println("Error generating rsa keypair in cs2a initialization")
		return err
	}
	rsaPublicKey := &rsaPrivateKey.PublicKey

	// generate the x509 certificate from the key pair
	template := gen_x509_template()
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, rsaPublicKey, rsaPrivateKey)
	if err != nil {
		log.Println("Error creating x509 certificate in cs2a initialization")
		log.Panicf("Error: %s\n")
	}
	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		log.Println("Error parsing x509 certificate in cs2a initialization")
		log.Panicf("Error: %s\n")
	}
	publicKey := cert.RawSubjectPublicKeyInfo

	// generate the fingerprint hash
	hash256 := sha256.New()
	hash256.Write(publicKey[:])
	fingerprintBin := hash256.Sum(nil)
	fingerprintHex := fmt.Sprintf("%x", fingerprintBin)

	// initialize the struct
	cs.id = "cs2a"
	cs.fingerprintBin = fingerprintBin
	cs.fingerprintHex = fingerprintHex
	cs.publicKey = publicKey
	cs.privateKey = *rsaPrivateKey
	cs.certificate = *cert

	return nil

}

// Generate a template x509 certificate
func gen_x509_template() *x509.Certificate {
	commonName := "thermal.telehash"
	organization := []string{"Acme Co"}
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 356)

	cert := new(x509.Certificate)
	cert.SerialNumber = big.NewInt(1)
	cert.Subject = pkix.Name{CommonName: commonName, Organization: organization}
	cert.NotBefore = notBefore
	cert.NotAfter = notAfter

	return cert
}

// fingerprint will return the csid and fingerprint for use in a 'parts' set
func (cs *cs2a) fingerprint() (string, string) {
	return cs.id, cs.fingerprintHex
}

// Stubs for the interface

func (cs *cs2a) encryptOpenPacket(packet []byte, remotePublicKey *[32]byte) (openPacketBody []byte, localLineSecret [32]byte, err error) {
	return openPacketBody, localLineSecret, err
}

func (cs *cs2a) decryptOpenPacket(openPacketBody []byte, remotePublicKey *[32]byte) (packet []byte, remoteLineSecret [32]byte, err error) {
	return packet, remoteLineSecret, err
}

func (cs *cs2a) generateLineEncryptionKey(localLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) (key [32]byte) {
	return key
}

func (cs *cs2a) generateLineDecryptionKey(remoteLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) (key [32]byte) {
	return key
}

func (cs *cs2a) encryptLinePacket(packet []byte, lineEncryptionKey *[32]byte) (linePacketBody []byte) {
	return linePacketBody
}

func (cs *cs2a) decryptLinePacket(linePacketBody []byte, lineDecryptionKey *[32]byte) (packet []byte) {
	return packet
}
