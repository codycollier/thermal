package thermal

import (
	"testing"
)

func TestHashingAccuracy(t *testing.T) {

	// known good parts and hashname taken from telehash docs
	knownParts := map[string]string{
		"2a": "bf6e23c6db99ed2d24b160e89a37c9cd183fb61afeca40c4bc378cf6e488bebe",
		"1a": "a5a741fa09b05baaead17fa9932e13cdafc7bcd39db1153fc6bbfe4614c063f3",
	}
	knownHash := "0b0137a6b38d00780686207b6f4b19e8731e68c6f76b435c85faf77100851451"

	testHash, err := generateHashname(knownParts)
	if err != nil {
		t.Fail()
	}
	if knownHash != testHash {
		t.Logf("knownHash: %s\n", knownHash)
		t.Logf(" testHash: %s\n", testHash)
		t.Log("Generated hash does not match expected hash")
		t.Fail()
	}
}

func TestHashnameGenerationBoundaries(t *testing.T) {

	// ensure empty parts causes error
	emptyParts := map[string]string{}
	_, err := generateHashname(emptyParts)
	if err == nil {
		t.Log("Empty parts should force error in generateHashname")
		t.Fail()
	}
}

func TestPartsExtraction(t *testing.T) {

	testCPack := make(cipherPack)
	cset := new(cs3a)
	cset.init()
	testCPack["cs3a"] = cset

	testparts, err := extractParts(&testCPack)

	// there should be no error
	if err != nil {
		t.Fail()
	}

	// there should be one entry per cipherSet
	if len(testparts) != 1 {
		t.Fail()
	}

	// the entry value should be a 64 len hex string
	if len(testparts["cs3a"]) != 64 {
		t.Logf("Bad fingerprint - testparts[cs3a]: %s\n", testparts["cs3a"])
		t.Fail()
	}

}

func TestPartsExtractionBounraries(t *testing.T) {

	var err error
	var cset cipherSet
	var cpack cipherPack

	// empty cpack should result in error
	cpack = make(cipherPack)
	_, err = extractParts(&cpack)
	if err == nil {
		t.Log("extractParts should return error for empty cpack")
		t.Fail()
	}

	// single non-initialized cipherSet should result in error
	cpack = make(cipherPack)
	cset = new(cs3a)
	cpack["cs3a"] = cset
	_, err = extractParts(&cpack)
	if err == nil {
		t.Log("extractParts should return error if it finds non-init cset")
		t.Fail()
	}

	// any non-initialized cipherSet should result in error
	cpack = make(cipherPack)

	cset = new(cs3a)
	cset.init()
	cpack["cs3a"] = cset

	cset = new(cs2a)
	cpack["cs2a"] = cset

	_, err = extractParts(&cpack)
	if err == nil {
		t.Log("extractParts should return error if it finds non-init cset")
		t.Fail()
	}

}
