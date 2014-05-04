package thermal

import (
	"testing"
)

func TestHashingAccuracy(t *testing.T) {
	knownParts := map[string]string{
		"2a": "bf6e23c6db99ed2d24b160e89a37c9cd183fb61afeca40c4bc378cf6e488bebe",
		"1a": "a5a741fa09b05baaead17fa9932e13cdafc7bcd39db1153fc6bbfe4614c063f3",
	}
	knownHash := "0b0137a6b38d00780686207b6f4b19e8731e68c6f76b435c85faf77100851451"
	testHash := generateHashname(knownParts)
	if knownHash != testHash {
		t.Logf("knownHash: %s\n", knownHash)
		t.Logf(" testHash: %s\n", testHash)
		t.Log("Generated hash does not match expected hash")
		t.Fail()
	}
}

func TestHashnameGenerationBoundaries(t *testing.T) {

	// ensure empty parts results in an empty hash
	emptyParts := map[string]string{}
	testHash := generateHashname(emptyParts)
	if testHash != "" {
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
		t.Logf("Bad fingerprint - testparts[cs3a]: %s", testparts["cs3a"])
		t.Fail()
	}

}

func TestPartsExtractionBounraries(t *testing.T) {

	var err error

	// empty cpack should result in empty parts
	emptyCPack := make(cipherPack)
	_, err = extractParts(&emptyCPack)
	if err == nil {
		t.Log("extractParts should return error for empty cpack")
		t.Fail()
	}

	// only non-initialized cipherSets should also result in empty parts
	noinitCPack := make(cipherPack)
	cset := new(cs3a)
	noinitCPack["cs3a"] = cset
	_, err = extractParts(&noinitCPack)
	if err == nil {
		t.Log("extractParts should return error if it finds non-init cset")
		t.Fail()
	}

}
