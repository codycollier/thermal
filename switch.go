package thermal

import (
	"fmt"
	"log"
)

// Switch provides an api to create, manage, and use a Telehash switch instance
type Switch struct {
	Hashname string

	// the internal backplane of the switch
	cpack     *cipherPack
	linestore *lineStore
	peerstore *peerStore
}

func (sw *Switch) String() string {
	return fmt.Sprintf("Switch(%s)", sw.Hashname)
}

// Initialize will setup all the internals of a switch instance
func (sw *Switch) Initialize(seedsPath, hintsPath string) error {

	// basic initialization
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	log.Println("Starting initialization of switch")

	// setup the cipher sets
	log.Println("Starting initialization of cipher sets")

	cpack := make(cipherPack)
	cpack["3a"] = new(cs3a)
	//cpack["2a"] = new(cs2a)

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
	sw.cpack = &cpack

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
	sw.Hashname = hashname
	log.Printf("switch hashname: %s", sw.Hashname)
	log.Println("Finished hashname creation")

	// Save identity file?
	idfile := fmt.Sprintf("./%s.id", hashname)
	writeIdentityFile(idfile, &cpack)

	// setup the line storage
	linestore := new(lineStore)
	linestore.start(sw)
	sw.linestore = linestore

	// setup the peer storage
	peerstore := new(peerStore)
	peerstore.start(sw)
	sw.peerstore = peerstore

	// load the seeds and hints files
	peerSeeds, err := loadPeersFile(seedsPath, "seed")
	if err != nil {
		log.Printf("Error loading seeds (%s)", seedsPath)
	}
	peerHints, err := loadPeersFile(hintsPath, "hint")
	if err != nil {
		log.Printf("Error loading hints (%s)", hintsPath)
	}
	_ = peerSeeds
	_ = peerHints

	// done
	log.Println("Finished initialization of switch")
	log.Println("Switch ready")
	return nil

}
