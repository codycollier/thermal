package thermal

import (
	"fmt"
)

// Swith provides an api to create, manage, and use a Telehash switch instance
type Switch struct{}

func (*Switch) Init() {
	fmt.Printf("Initializing switch\n")
	cstest := new(cs3a)
	cstest.init()
	fmt.Printf("cstest: %s\n", cstest)
}
