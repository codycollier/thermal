package thermal

import (
	"bytes"
	"encoding/gob"
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

	if peersFile != "" {
		peersJson, err := ioutil.ReadFile(peersFile)
		if err != nil {
			log.Printf("Error reading file (err: %s)", err)
			return peers, err
		}
		peers = loadPeersFromJson(peersJson)
	}

	return peers, nil
}

// writeIdentityFile will save the local switch identity information to file
func writeIdentityFile(idFileName string, cpack *cipherPack) error {

	gob.Register(cipherPack{})
	gob.Register(&cs3a{})
	gob.Register(&cs2a{})

	var encodedFileData bytes.Buffer
	enc := gob.NewEncoder(&encodedFileData)

	err := enc.Encode(cpack)
	if err != nil {
		log.Printf("Error encoding data (err: %s)", err)
		return err
	}

	err = ioutil.WriteFile(idFileName, encodedFileData.Bytes(), 0644)
	if err != nil {
		log.Printf("Error writing file (file: %s) (err: %s)", idFileName, err)
		return err
	}

	return nil
}

// readIdentityFile will read the local switch identity information from file
func readIdentityFile(idFileName string, cpack *cipherPack) error {

	gob.Register(cipherPack{})
	gob.Register(&cs3a{})
	gob.Register(&cs2a{})

	//var encodedFileData bytes.Buffer
	var encodedFileData []byte

	encodedFileData, err := ioutil.ReadFile(idFileName)
	if err != nil {
		log.Printf("Error reading file (file: %s) (err: %s)", idFileName, err)
		return err
	}

	dec := gob.NewDecoder(bytes.NewReader(encodedFileData))
	err = dec.Decode(cpack)
	if err != nil {
		log.Printf("Error decoding data (err: %s)", err)
		return err
	}

	return nil
}
