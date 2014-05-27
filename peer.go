package thermal

/*
{
"46fe53c258bbc1984fb5ab02ca1494eccdd54e9688dbbc2c882c8713f1cc4cf3":{
    "admin":"http://github.com/quartzjer",
    "paths":[{"type":"ipv4","ip": "127.0.0.1","port": 42424},{"type":"http","http":"http://127.0.0.1"}],
    "keys":{
      "1a":"z6yCAC7r5XIr6C4xdxeX7RlSmGu9Xe73L1gv8qecm4/UEZAKR5iCxA==",
      "2a":"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnDQ/EdMwXn3nAGaEH3bM37xbG71M41iQTnE56xh+RS8kvjAaEG3mxqcezEFyLTuhb8oraoQeHvD8mmCdm+NNpuYUgx3SmnwGO91JsVnVHi94kL5P9UzT501k43nJq+Lnjx5FamFyDDVulAGiOuw4HQHqBuiGsjqQzRO7CclQtlBNewPQUrwoVG7K60+8EIpNuD6opyC6fH1XYNtx10G8hyN1bEyRN+9xsgW3I8Yw8sbPjFhuZGfM0nlgevdG4n+cJaG0fVdag1tx08JiWDlYm3wUWCivLeQTOLKrkVULnPw06YxvWdUURg742avZqMKhZTGsHJgHJir3Tfw9kk0eFwIDAQAB"
    },
    "parts":{
      "1a":"b5a96d25802b3600ea99774138a650d5d1fa1f3cf3cb10ae8f1c58a527d85086",
      "2a":"40a344de8c6e93282d085c577583266e18ed23182d64e382b7e31e05fec57d67"
    }
  }
}
*/

type path struct {
	pathtype string
	ip       string
	port     string
	http     string
}

type peerSwitch struct {
	hashname string
	admin    string
	paths    []path
	keys     map[string]string
	parts    map[string]string
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
