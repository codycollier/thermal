package thermal

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestSwitchInit(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	s := new(Switch)
	s.Init()
}
