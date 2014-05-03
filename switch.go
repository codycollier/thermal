package thermal

import (
	"log"
)

// Switch provides an api to create, manage, and use a Telehash switch instance
type Switch struct{}

func (*Switch) Init() {

	// basic initialization
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	log.Println("Starting initialization of switch")

	// setup the cipher sets
	log.Println("Starting initialization of cipher sets")
	cpack := make(cipherPack)

	cpack["cs2a"] = new(cs2a)
	cpack["cs3a"] = new(cs3a)

	for csid, cset := range cpack {

		log.Printf("init %s...\n", csid)
		cset.init()

		id, fingerprint := cset.fingerprint()
		log.Printf("       csid: %s\n", id)
		log.Printf("fingerprint: %s\n", fingerprint)

	}
	log.Println("Finished initialization of cipher sets")

	//
	log.Println("Finished initialization of switch")
	log.Println("Switch ready")

}
