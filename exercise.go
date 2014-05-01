package main

import (
	"fmt"
	"github.com/codycollier/thermal/thermal"
	"time"
)

// exercise will exercise and test the thermal library
func main() {
	start := time.Now()
	fmt.Printf("start: %s\n", start)
	s := new(thermal.Switch)
	fmt.Printf("Switch instance created: %s\n", *s)
	fmt.Printf("stop: %s\n", time.Now())
}
