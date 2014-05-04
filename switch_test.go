package thermal

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestSwitchInit(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping switch init test in short mode")
	}

	log.SetOutput(ioutil.Discard)
	s := new(Switch)
	err := s.Init()
	if err != nil {
		t.Fail()
	}
}
