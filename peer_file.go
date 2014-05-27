package thermal

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type peerLoader map[string]struct {
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

// loadPeersFromJson returns peers unpacked from json
func loadPeersFromJson(peerjson []byte) []peerSwitch {

	var peers []peerSwitch
	var peerloader peerLoader

	// unpack each 'seed' in to a peer struct
	json.Unmarshal(peerjson, &peerloader)
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

// loadPeersFile returns peers loaded from a seeds.json or hints.json file
func loadPeersFile(peersFile, peerType string) ([]peerSwitch, error) {

	var peers []peerSwitch
	var err error

	peersJson, err := ioutil.ReadFile(peersFile)
	if err != nil {
		log.Printf("Error reading file (%s) (err: %s)", peersFile, err)
		return peers, err
	}
	peers = loadPeersFromJson(peersJson)

	return peers, nil
}
