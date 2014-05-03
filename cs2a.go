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
	id          string
	fingerprint string
	publicKey   []byte
	privateKey  rsa.PrivateKey
	certificate x509.Certificate
}

// init will generate a key pair and initialize the cipher set
func (cs *cs2a) init() {

	// generate the rsa-2048 key pair
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println("Error generating rsa keypair in cs2a initialization")
		log.Panicf("Error: %s\n", err)
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
	fingerprint_bytes := hash256.Sum(nil)
	fingerprint_hex := fmt.Sprintf("%x", fingerprint_bytes)

	// initialize the struct
	cs.id = "cs2a"
	cs.fingerprint = fingerprint_hex
	cs.publicKey = publicKey
	cs.privateKey = *rsaPrivateKey
	cs.certificate = *cert

}

func (cs *cs2a) String() string {
	return fmt.Sprintf("%s: %x", cs.id, cs.fingerprint)
}

// csid will return the id of the cipher set
func (cs *cs2a) csid() string {
	return cs.id
}

// parts will return the telehash defined 'parts' for the cipherset
func (cs *cs2a) parts() (string, string) {
	return cs.id, cs.fingerprint
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
