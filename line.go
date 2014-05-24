package thermal

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"
)

//-----------------------------------------------------------------------------
// Router
// The router helps decouple lines from channels and paths.
//-----------------------------------------------------------------------------

// routeToLine will route packets to their appropriate line
func routeToLine() {
}

//-----------------------------------------------------------------------------
// Line Store
// The line store is a service which manages the mapping of remote switches
// to lines.
//-----------------------------------------------------------------------------

// The storeRequest represents requests for lines
type storeRequest struct {
	hashname string
	response chan *lineSession
}

// The lineStore holds and manages the lines
type lineStore struct {
	lineMap map[string]*lineSession
	request chan *storeRequest
}

// The service func listens and handles incoming requests
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

// start will setup the listener to service requests
func (store *lineStore) start() {
	go store.service()
}

//-----------------------------------------------------------------------------
// Line(s)
// Representation and logic for telehash lines
//-----------------------------------------------------------------------------

// A lineHalf is one half (local or remote) of a Line
type lineHalf struct {
	id     string
	at     int64
	secret [32]byte
}

// A lineSession represents a telehash Line between two switches
type lineSession struct {
	cset            cipherSet
	remoteHashname  string
	remotePublicKey *[32]byte

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

// service will listen and respond to open/send/recv messages
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

// start will setup the line listener
func (line *lineSession) start(remoteHashname string) {

	// set some line vars
	// todo - where to get the cset and remotePublickKey?
	//		get the public key for a remote hashname from...
	//		maybe pass them in with the storeRequest()
	line.cset = cset
	line.remoteHashname = remoteHashname
	line.remotePublicKey = remotePublicKey

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

// newLocalLine will...
func (line *lineSession) newLocalLine() {

	line.local.id = generateLineId()
	line.local.at = generateLineAt()

	// todo
	// json := make the json (to, from(parts), at, localLineId)
	// to == remoteHashname
	// parts will be retrieved over in the openMaker()?
	// lin.local.id
	json := "{}"
	body := line.cset.pubKey()[:]
	packet, err := encodePacket(json, body)
	openPacketBody, localLineSecret, err := line.cset.encryptOpenPacketBody(packet, line.remotePublicKey)
	if err != nil {
		log.Printf("Error encrypting open packet body (err: %s)", err)
		return
	}

	line.local.secret = localLineSecret

	// todo
	// the head needs to be a single byte, representing the csid
	// the encodePacket() function is not setup to handle that as-is
	openPacketJson := ""
	openPacket, err := encodePacket(openPacketJson, openPacketBody)
	if err != nil {
		log.Printf("Error encoding open packet (err: %s)", err)
		return
	}
	log.Println(openPacket)

	// todo
	// return or send
}

// newRemoteLine will...
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
