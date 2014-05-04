package thermal

import (
	"testing"
)

func TestCs3aInit(t *testing.T) {

	cset := new(cs3a)
	err := cset.init()
	if err != nil {
		t.Log("Error in cs3a initialization")
		t.Logf("err: %s\n", err)
		t.Fail()
	}

}

func TestCs3aCsid(t *testing.T) {
	cset := new(cs3a)
	cset.init()
	csid := cset.csid()
	if csid != "cs3a" {
		t.Fail()
	}
}

func TestCs3aFingerprint(t *testing.T) {
	cset := new(cs3a)
	cset.init()
	csid, fingerprint := cset.fingerprint()
	if csid != "cs3a" {
		t.Fail()
	}
	if len(fingerprint) != 64 {
		t.Fail()
	}

}
