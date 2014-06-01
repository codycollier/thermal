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

	idFile := ""
	seedsFile := "./test/seeds-test.json"
	hintsFile := ""

	err := s.Initialize(idFile, seedsFile, hintsFile)
	if err != nil {
		t.Fail()
	}
}
