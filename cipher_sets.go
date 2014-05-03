package thermal

// The cipherSet interface defines behavior expected from a cipher set plugin
type cipherSet interface {
	init()
	csid() string
	parts() (string, string)
}

// A cipherPack holds a group of cipher sets
type cipherPack map[string]cipherSet
