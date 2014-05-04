package thermal

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestSwitchInit(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	s := new(Switch)
	err := s.Init()
	if err != nil {
		t.Fail()
	}
}
