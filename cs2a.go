package thermal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"
)

// cs2a is an implementation of the cipher set 2a
type cs2a struct {
	id             [1]byte
	publicKey      [32]byte
	privateKey     rsa.PrivateKey
	fingerprintBin []byte

	certificate x509.Certificate
}

// init will generate a key pair and setup the cipher set
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

	cs.populate(publicKey, rsaPrivateKey)
	return nil
}

// populate takes a generated or loaded key pair and sets up the cipher set
func (cs *cs2a) populate(publicKey []byte, rsaPrivateKey *rsa.PrivateKey) {

	// generate the fingerprint hash
	hash256 := sha256.New()
	hash256.Write(publicKey[:])
	fingerprintBin := hash256.Sum(nil)

	// initialize the struct
	csid_byte, _ := hex.DecodeString("2a")
	copy(cs.id[:], csid_byte[:])
	cs.fingerprintBin = fingerprintBin
	copy(cs.publicKey[:], publicKey)
	cs.privateKey = *rsaPrivateKey
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

func (cs *cs2a) String() string {
	return fmt.Sprintf("%x: %x", cs.id, cs.fingerprint)
}

func (cs *cs2a) csid() [1]byte {
	return cs.id
}

// fingerprint will return the csid and fingerprint for use in a 'parts' set
func (cs *cs2a) fingerprint() (string, string) {
	return fmt.Sprintf("%x", cs.id), fmt.Sprintf("%x", cs.fingerprintBin)
}

// pubKey returns the cipher set public key
func (cs *cs2a) pubKey() *[32]byte {
	return &cs.publicKey
}

// Stubs for the interface

func (cs *cs2a) encryptOpenPacketBody(packet []byte, remotePublicKey *[32]byte) (openPacketBody []byte, localLineSecret [32]byte, err error) {
	return openPacketBody, localLineSecret, err
}

func (cs *cs2a) decryptOpenPacketBody(openPacketBody []byte, remotePublicKey *[32]byte) (packet []byte, remoteLineSecret [32]byte, err error) {
	return packet, remoteLineSecret, err
}

func (cs *cs2a) generateLineEncryptionKey(localLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) (key [32]byte) {
	return key
}

func (cs *cs2a) generateLineDecryptionKey(remoteLineSecret *[32]byte, localLineId, remoteLineId *[16]byte) (key [32]byte) {
	return key
}

func (cs *cs2a) encryptLinePacketBody(packet []byte, lineEncryptionKey *[32]byte) (linePacketBody []byte, err error) {
	return linePacketBody, nil
}

func (cs *cs2a) decryptLinePacketBody(linePacketBody []byte, lineDecryptionKey *[32]byte) (packet []byte, err error) {
	return packet, nil
}
