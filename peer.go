package thermal

import (
	"fmt"
	"log"
)

// A path is network transport information
type path struct {
	pathtype string
	ip       string
	port     int64
	http     string
}

// A peerSwitch is a representation of a remote telehash switch
type peerSwitch struct {
	hashname string
	admin    string
	paths    []path
	keys     map[string]string
	parts    map[string]string
}

func (p *peerSwitch) String() string {
	return fmt.Sprintf("peerSwitch(%s)", p.hashname)
}

// peerStoreRequests are used to communicate with the store
type peerStoreRequest struct {
	hashname string
	peerdata *peerSwitch
	response chan *peerSwitch
}

// The peerStore holds and manages information about other telehash switches
type peerStore struct {
	sw       *Switch
	peerMap  map[string]peerSwitch
	requests chan *peerStoreRequest
}

// service is the main loop for handling incoming messages to the store
func (ps *peerStore) service() {
	log.Println("Starting peerstore service")
	for {

		select {

		case request := <-ps.requests:

			log.Println("peerstore: accepted a request")

			if request.peerdata != nil {

				// Update or Create a peer entry
				valid := validatePeer(request.peerdata)
				if !valid {
					request.response <- nil
				}
				ps.peerMap[request.hashname] = *request.peerdata
				log.Println("peerstore: sending a response")
				request.response <- request.peerdata

			} else {

				// Get an existing peer entry or create a new, empty peer entry
				peer, exists := ps.peerMap[request.hashname]
				if !exists {
					var peer peerSwitch
					peer.hashname = request.hashname
					ps.peerMap[request.hashname] = peer
				}
				log.Println("peerstore: sending a response")
				request.response <- &peer
			}

		default:
			//
			//log.Println("default...")
		}
	}
}

// start will setup the peerStore and start the listening service
func (ps *peerStore) start(sw *Switch) {
	ps.sw = sw
	ps.peerMap = make(map[string]peerSwitch)
	ps.requests = make(chan *peerStoreRequest)
	go ps.service()
}

// TODO - validate incoming peerSwitch data
func validatePeer(p *peerSwitch) bool {
	return true
}
