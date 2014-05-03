package thermal

import (
	"log"
)

// Switch provides an api to create, manage, and use a Telehash switch instance
type Switch struct{}

func (*Switch) Init() {

	// basic initialization
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	log.Printf("Initializing switch\n")

	// setup the cipher sets
	cpack := make(cipherPack)

	cpack["cs2a"] = new(cs2a)
	cpack["cs3a"] = new(cs3a)

	for csid, cset := range cpack {

		cset.init()

		id := cset.csid()
		log.Printf(" csid: %s :: %s\n", id, csid)

		id, fingerprint := cset.parts()
		log.Printf("parts: %s :: %s\n", id, fingerprint)

	}

}
