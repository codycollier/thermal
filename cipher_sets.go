package thermal

// cipherSetPlugin is the base type for all cipher set types
type cipherSetPlugin struct {
	id          string
	fingerprint []byte
	private_key [32]byte
	public_key  [32]byte
}

// The cypherSet interface defines behavior expected from a cipher set plugin
type cipherSet interface{}

// The cipherPack holds all the cipher sets for a switch instance
type cipherPack map[string]cipherSet
