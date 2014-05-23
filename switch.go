package thermal

import (
	"log"
)

// Switch provides an api to create, manage, and use a Telehash switch instance
type Switch struct {
	hashname string
}

func (s *Switch) Initialize() error {

	// basic initialization
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	log.Println("Starting initialization of switch")

	// setup the cipher sets
	log.Println("Starting initialization of cipher sets")

	cpack := make(cipherPack)
	//cpack["cs2a"] = new(cs2a)
	cpack["cs3a"] = new(cs3a)

	for csid, cset := range cpack {

		log.Printf("initialize %s...\n", csid)

		err := cset.initialize()
		if err != nil {
			log.Printf("Error in initialization of cset. csid:%s  err: %s", csid, err)
			return err
		}

		id, fingerprint := cset.fingerprint()
		log.Printf("       csid: %s\n", id)
		log.Printf("fingerprint: %s\n", fingerprint)

	}
	log.Println("Finished initialization of cipher sets")

	// build the hashname for the switch instance
	log.Println("Starting hashname creation")
	parts, err := extractParts(&cpack)
	if err != nil {
		return err
	}
	hashname, err := generateHashname(parts)
	if err != nil {
		return err
	}
	s.hashname = hashname
	log.Printf("switch hashname: %s", s.hashname)
	log.Println("Finished hashname creation")

	// setup the line storage
	linestore := new(lineStore)
	linestore.start()

	// end
	log.Println("Finished initialization of switch")
	log.Println("Switch ready")
	return nil

}
