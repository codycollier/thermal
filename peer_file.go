package thermal

import (
	"encoding/json"
)

/*

// seeds from docs:
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

// seeds.json from telehash-c project:
{
  "dca549c98b94197e79dfa2a7e4ad2e74fb144cda3f4710d4f40e2c75d975272e": {
    "paths": [
      {
        "type": "http",
        "http": "http://192.168.0.36:42424"
      },
      {
        "type": "ipv4",
        "ip": "127.0.0.1",
        "port": 42424
      },
      {
        "type": "ipv6",
        "ip": "fe80::bae8:56ff:fe43:3de4",
        "port": 42424
      }
    ],
    "parts": {
      "3a": "f0d2bfc8590a7e0016ce85dbf0f8f1883fb4f3dcc4701eab12ef83f972a2b87f",
      "2a": "0cb4f6137a745f1af2d31707550c03b99083180f6e69ec37918c220ecfa2972f",
      "1a": "b5a96d25802b3600ea99774138a650d5d1fa1f3cf3cb10ae8f1c58a527d85086"
    },
    "keys": {
      "3a": "MC5dfSfrAVCSugX75JbgVWtvCbxPqwLDUkc9TcS/qxE=",
      "2a": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqr12tXnpn707llkZfEcspB/D6KTcZM765+SnI5Z8JWkjc0Mrz9qZBB2YFLr2NmgCx0oLfSetmuHBNTT54sIAxQ/vxyykcMNGsSFg4WKhbsQXSrX4qChbhpIqMJkKa4mYZIb6qONA76G5/431u4+1sBRvfY0ewHChqGh0oThcaa50nT68f8ohIs1iUFm+SL8L9UL/oKN3Yg6drBYwpJi2Ex5Idyu4YQJwZ9sAQU49Pfs+LqhkHOascTmaa3+kTyTnp2iJ9wEuPg+AR3PJwxXnwYoWbH+Wr8gY6iLe0FQe8jXk6eLw9mqOhUcah8338MC83zSQcZriGVMq8qaQz0L9nwIDAQAB",
      "1a": "z6yCAC7r5XIr6C4xdxeX7RlSmGu9Xe73L1gv8qecm4/UEZAKR5iCxA=="
    }
  }
}


*/

type peerLoader1 map[string]peerLoader2

type peerLoader2 struct {
	Admin string
	Paths []pathLoader
	Parts map[string]string
	Keys  map[string]string
}

type pathLoader struct {
	Pathtype string `json:"Type"`
	Ip       string
	Port     int64
	Http     string
}

func loadPeersFromString(peerlist string) []peerSwitch {

	var peers []peerSwitch
	var peerloader peerLoader1

	// unpack each 'seed' in to a peer struct
	json.Unmarshal([]byte(peerlist), &peerloader)
	for hashname, pl2 := range peerloader {

		// map pathLoader(s) to path(s)
		paths := make([]path, 0)
		for _, p := range pl2.Paths {
			paths = append(paths, path{p.Pathtype, p.Ip, p.Port, p.Http})
		}

		peers = append(peers, peerSwitch{hashname, pl2.Admin, paths, pl2.Keys, pl2.Parts})
	}

	return peers
}
