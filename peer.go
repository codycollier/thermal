package thermal

import (
	"fmt"
)

type path struct {
	pathtype string
	ip       string
	port     int64
	http     string
}

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

type peerStoreRequest struct {
	hashname string
	peerdata *peerSwitch
	response chan *peerSwitch
}

type peerStore struct {
	sw       *Switch
	peerMap  map[string]peerSwitch
	requests chan *peerStoreRequest
}

func (ps *peerStore) service() {
	for {
		select {

		case request := <-ps.requests:
			if request.peerdata != nil {

				// If there is an incoming peer value, then update the store
				valid := validatePeer(request.peerdata)
				if !valid {
					request.response <- nil
				}
				ps.peerMap[request.hashname] = *request.peerdata
				request.response <- request.peerdata

			} else {

				// Get or create peer entry
				peer, exists := ps.peerMap[request.hashname]
				if !exists {
					var peer peerSwitch
					peer.hashname = request.hashname
					ps.peerMap[request.hashname] = peer
				}
				request.response <- &peer
			}

		default:
			//
		}
	}
}

func (ps *peerStore) start(sw *Switch) {
	ps.sw = sw
	go ps.service()
}

// TODO - validate?
func validatePeer(p *peerSwitch) bool {
	return true
}
