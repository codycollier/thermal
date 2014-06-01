package main

import (
	"github.com/codycollier/thermal"
	"log"
	"time"
)

func startPing() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	start := time.Now()
	log.Println("Starting exercise")

	idFile := "./501ae97f1b5fcf49af34c9ce53279574ee87fffcc667b659b5a8163d7e48441f.id"
	seedsFile := "./seeds-for-ping.json"
	hintsFile := ""

	s := new(thermal.Switch)
	s.Initialize(idFile, seedsFile, hintsFile)

	log.Printf("Switch instance created: %s", s)
	log.Printf("switch hashname: %s", s.Hashname)
	log.Printf("Finished exercise (%s)\n", time.Since(start))
}

func main() {
	startPing()
}
