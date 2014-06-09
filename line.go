package thermal

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

// routeToLine will route packets to their appropriate line
func routeToLine() {
}

//-----------------------------------------------------------------------------
// Line Store
// The line store is a service which manages access to lines
//-----------------------------------------------------------------------------

type lineStoreRequest struct {
	peer     *peerSwitch
	response chan *lineSession
}

type lineStore struct {
	sw      *Switch
	lineMap map[string]*lineSession
	request chan *lineStoreRequest
}

func (store *lineStore) service() {
	for {
		select {

		case request := <-store.request:
			// Get or Create a line for the hashname
			line, exists := store.lineMap[request.peer.hashname]
			if !exists {
				line := new(lineSession)
				line.start(store.sw, request.peer)
				store.lineMap[request.peer.hashname] = line
			}
			request.response <- line

		default:
		}
	}
}

func (store *lineStore) start(sw *Switch) {
	store.sw = sw
	go store.service()
}

//-----------------------------------------------------------------------------
// Line(s)
// Representation and logic for telehash lines
//-----------------------------------------------------------------------------

type lineHalf struct {
	id     string
	at     int64
	secret [32]byte
}

type lineSession struct {
	csid            string
	cset            cipherSet
	remotePublicKey [32]byte

	local         lineHalf
	remote        lineHalf
	encryptionKey [32]byte
	decryptionKey [32]byte

	ready      bool
	openLocal  chan bool
	openRemote chan []byte
	send       chan decodedPacket
	recv       chan []byte
}

func (line *lineSession) service() {

	for {

		select {

		case msg := <-line.openLocal:
			log.Println("openLocal...")
			log.Printf("msg: %s", msg)
			//line.newLocalLine()

		case msg := <-line.openRemote:
			log.Println("openRemote...")
			log.Printf("msg: %s", msg)
			//line.openRemote()

		default:
			//
		}

		if line.ready {
			select {

			case msg := <-line.send:
				log.Println("line send request...")
				log.Printf("msg: %s", msg)

			case msg := <-line.recv:
				log.Println("line recv request...")
				log.Printf("msg: %s", msg)

			default:

			}
		}

	}
}

func (line *lineSession) start(sw *Switch, peer *peerSwitch) {

	// determine best cipher set which both switches support
	// todo - getBestCSMatch(sw, peer)
	agreedCSID := "3a"
	line.csid = agreedCSID
	// todo - better way to dereference
	cp := *sw.cpack
	line.cset = cp[line.csid]

	// set the remotePublicKey for later use (decode from base64)
	peerKey, _ := base64.StdEncoding.DecodeString(peer.keys[line.csid])
	copy(line.remotePublicKey[:], peerKey[:])

	// setup the channels
	line.openLocal = make(chan bool)
	line.openRemote = make(chan []byte)
	line.send = make(chan decodedPacket)
	line.recv = make(chan []byte)

	// queue up a message to begin the open process
	line.ready = false
	line.openLocal <- true

	go line.service()
}

func (line *lineSession) newLocalLine() {

	line.local.id = generateLineId()
	line.local.at = generateLineAt()

	// todo
	// json := make the json (to, from(parts), at, localLineId)
	// to == remoteHashname
	// parts will be retrieved over in the openMaker()?
	// line.local.id
	cset := line.cset
	jsonHead := []byte("{}")
	body := cset.pubKey()[:]
	packet, err := encodePacket(jsonHead, body)
	openPacketBody, localLineSecret, err := cset.encryptOpenPacketBody(packet, &line.remotePublicKey)
	if err != nil {
		log.Printf("Error encrypting open packet body (err: %s)", err)
		return
	}

	line.local.secret = localLineSecret

	// todo
	// the head needs to be a single byte, representing the csid
	// the encodePacket() function is not setup to handle that as-is
	openPacketHead := []byte("")
	openPacket, err := encodePacket(openPacketHead, openPacketBody)
	if err != nil {
		log.Printf("Error encoding open packet (err: %s)", err)
		return
	}
	log.Println(openPacket)

	// todo
	// return or send
}

func (line *lineSession) newRemoteLine() {
}

// generateLineId returns a random 16 char hex encoded string
func generateLineId() string {
	var idBin [8]byte
	rand.Reader.Read(idBin[:])
	idHex := fmt.Sprintf("%x", idBin)
	return idHex
}

// generateLineAt returds an integer timestamp suitable for a line at
func generateLineAt() int64 {
	return time.Now().Unix()
}
