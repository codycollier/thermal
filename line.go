package thermal

import (
	"crypto/rand"
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
	hashname string
	// what about a pointer to a peer entry?  Or do we look that up later?
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
			line, exists := store.lineMap[request.hashname]
			if !exists {
				line := new(lineSession)
				line.start(request.hashname)
				store.lineMap[request.hashname] = line
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
	remoteHashname  string
	remotePublicKey *[32]byte
	cset            cipherSet

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

		case <-line.openLocal:
			//line.newLocalLine()

		case something := <-line.openRemote:
			//line.openRemote()
			log.Println(something)

		default:
			//
		}

		if line.ready {
			select {
			case something := <-line.send:
				log.Println(something)
			case something := <-line.recv:
				log.Println(something)
			default:
			}
		}

	}
}

func (line *lineSession) start(remoteHashname string) {

	// set some line vars
	// todo - where to get the cset and remotePublickKey?
	//		get the public key for a remote hashname from...
	//		maybe pass them in with the storeRequest()
	line.remoteHashname = remoteHashname
	//line.cset = cset
	//line.remotePublicKey = remotePublicKey

	// setup the channels
	line.openLocal = make(chan bool)
	line.openRemote = make(chan []byte)
	line.send = make(chan decodedPacket)
	line.recv = make(chan []byte)

	// initialization has not completed yet
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
	jsonHead := []byte("{}")
	body := line.cset.pubKey()[:]
	packet, err := encodePacket(jsonHead, body)
	openPacketBody, localLineSecret, err := line.cset.encryptOpenPacketBody(packet, line.remotePublicKey)
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
