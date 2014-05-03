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

	s := new(thermal.Switch)
	s.Init()
	log.Printf("Switch instance created: %s\n", *s)

	log.Printf("Finished exercise (%s)\n", time.Since(start))
}

func main() {
	exercise()
}
