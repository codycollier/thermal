package thermal

import (
	"log"
)

// Swith provides an api to create, manage, and use a Telehash switch instance
type Switch struct{}

func (*Switch) Init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	log.Printf("Initializing switch\n")

	cstest := new(cs3a)
	cstest.init()

	log.Printf("cs3a fingerprint: %x\n", cstest.fingerprint)
}
