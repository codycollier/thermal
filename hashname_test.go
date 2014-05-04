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
	extractParts(&testCPack)

}

func TestPartsExtractionBounraries(t *testing.T) {

	var testparts map[string]string

	// totally empty cpack
	emptyCPack := make(cipherPack)
	testparts = extractParts(&emptyCPack)
	if len(testparts) != 0 {
		t.Fail()
	}

	// cpack with a non-initialized cipherSet
	noinitCPack := make(cipherPack)
	cset := new(cs3a)
	noinitCPack["cs3a"] = cset
	extractParts(&noinitCPack)
	if len(testparts) != 0 {
		t.Fail()
	}

}
