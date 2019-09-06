package pkpass

// SignatureFilename is an alias for `signature`.
const SignatureFilename = "signature"

// Signer holds method working with certificates.
type Signer interface {
	Sign(data []byte) (*File, error)
}
