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
func (sw *Switch) Initialize(idFile, seedsPath, hintsPath string) error {

	var err error
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	log.Println("Starting initialization of switch")

	if idFile != "" {
		// Read in a pre-existing identity / cipherPack
		log.Println("Reading in pre-existing cipher pack")
		cpack := make(cipherPack)
		sw.cpack = &cpack
		readIdentityFile(idFile, sw.cpack)
		for csid, cset := range *sw.cpack {
			_, fingerprint := cset.fingerprint()
			pubkey := cset.pubKeyStr()
			log.Printf("       csid: %s\n", csid)
			log.Printf("fingerprint: %s\n", fingerprint)
			log.Printf(" public_key: %s\n", pubkey)
		}

	} else {
		// Generate a new cipherPack
		log.Println("Starting initialization of new cipher pack")
		err = sw.newCipherPack()
		if err != nil {
			return err
		}
	}

	log.Println("Generating hashname from cipher pack")
	err = sw.newHashname()
	if err != nil {
		return err
	}

	idfile := fmt.Sprintf("./%s.id", sw.Hashname)
	writeIdentityFile(idfile, sw.cpack)

	sw.initializeStores()

	if seedsPath != "" {
		err = sw.loadPeers(seedsPath, "seeds")

		if err != nil {
			log.Printf("Error loading seeds (%s)", seedsPath)
		}
	}
	if hintsPath != "" {
		err = sw.loadPeers(hintsPath, "hints")
		if err != nil {
			log.Printf("Error loading hints (%s)", hintsPath)
		}
	}

	log.Println("Finished initialization of switch")
	log.Println("Switch ready")
	return nil

}

// newCipherPack generates and sets a new cipherPack
func (sw *Switch) newCipherPack() error {

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
	return nil
}

// newHashname generates and sets a new local hashname
func (sw *Switch) newHashname() error {

	log.Println("Starting hashname creation")
	parts, err := extractParts(sw.cpack)
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
	return nil
}

// initializeStores sets up the backplane storage
func (sw *Switch) initializeStores() {

	linestore := new(lineStore)
	linestore.start(sw)
	sw.linestore = linestore

	peerstore := new(peerStore)
	peerstore.start(sw)
	sw.peerstore = peerstore
}

// loadPeers populates the peerstore with pre-existing seeds or hints
func (sw *Switch) loadPeers(peersFile, peersType string) error {

	peers, err := loadPeersFile(peersFile, peersType)
	log.Printf("Loaded peers of type %s", peersType)
	log.Printf("Loaded peers %s", peers)
	for peer := range peers {
		log.Printf("peer: %s", peer)
	}

	if err != nil {
		return err
	}

	_ = peers

	return nil

}
