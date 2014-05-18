package thermal

import (
	"crypto/sha256"
	"fmt"
	"sort"
)

// A hashname is a 64 byte hex string
type hashname string

// extractParts will pull out the hashname parts from a cipherPack
func extractParts(cpack *cipherPack) (map[string]string, error) {

	if len(*cpack) == 0 {
		return nil, fmt.Errorf("cipherPack is empty")
	}

	parts := make(map[string]string)

	for _, cset := range *cpack {
		csid, fingerprint := cset.fingerprint()

		if fingerprint == "" {
			return nil, fmt.Errorf("fingerprint for %s is empty", csid)
		}

		parts[csid] = fingerprint
	}

	return parts, nil
}

// generateHashname will create a switch hashname from the given parts
func generateHashname(parts map[string]string) (hashname, error) {

	if len(parts) == 0 {
		return "", fmt.Errorf("parts map is empty")
	}

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
		hash.Write(append(rollup, []byte(csid)...))
		rollup = hash.Sum(nil)

		hash.Reset()
		fingerprint := parts[csid]
		hash.Write(append(rollup, []byte(fingerprint)...))
		rollup = hash.Sum(nil)

	}

	hashnameBin := rollup
	hashnameHex := hashname(fmt.Sprintf("%x", hashnameBin))

	return hashnameHex, nil
}
