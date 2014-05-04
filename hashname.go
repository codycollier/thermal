package thermal

import (
	"crypto/sha256"
	"fmt"
	"sort"
)

// extractParts will pull out the hashname parts from a cipherPack
func extractParts(cpack *cipherPack) map[string]string {

	parts := make(map[string]string)

	for _, cset := range *cpack {
		csid, fingerprint := cset.fingerprint()
		parts[csid] = fingerprint
	}

	return parts
}

// generateHashname will create a switch hashname from the given parts
func generateHashname(parts map[string]string) string {

	var keys []string
	rollup := make([]byte, 0)
	hash := sha256.New()

	// The rollup hash is specified to be done in order of csid
	for key := range parts {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Generate the hash in rollup fashion
	for _, csid := range keys {

		hash.Reset()
		hash.Write([]byte(csid))
		rollup = hash.Sum(rollup)

		fingerprint := parts[csid]
		hash.Reset()
		hash.Write([]byte(fingerprint))
		rollup = hash.Sum(rollup)

	}

	hashnameBin := rollup
	hashnameHex := fmt.Sprintf("%x", hashnameBin)

	return hashnameHex
}
