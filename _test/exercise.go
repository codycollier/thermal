package main

import (
	"github.com/codycollier/thermal"
	"log"
	"time"
)

// exercise will exercise and test the thermal library
func exercise() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	start := time.Now()
	log.Println("Starting exercise")

	idFile := ""
	seedsFile := "./seeds-test.json"
	hintsFile := ""

	s := new(thermal.Switch)
	s.Initialize(idFile, seedsFile, hintsFile)

	log.Printf("Switch instance created: %s", s)
	log.Printf("switch hashname: %s", s.Hashname)
	log.Printf("Finished exercise (%s)\n", time.Since(start))
}

func main() {
	exercise()
	//panic("")
}
