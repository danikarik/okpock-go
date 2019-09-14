package pkpass

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

// Manifest holds file name and SHA1 hash as a map.
type Manifest map[string]string

// HashFile generates SHA1 hash of file content.
func HashFile(content []byte) (string, error) {
	h := sha1.New()
	_, err := h.Write(content)
	if err != nil {
		return "", err
	}
	sum := h.Sum(nil)
	hash := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(hash, sum)
	return string(hash), nil
}

// CreateManifest generates File called `manifest.json`.
func CreateManifest(files ...File) (*File, error) {
	if len(files) == 0 {
		return nil, ErrEmptyFolder
	}

	manifest := make(Manifest)
	for _, file := range files {
		hash, err := HashFile(file.Data)
		if err != nil {
			return nil, err
		}
		manifest[file.Name] = hash
	}

	data, err := json.Marshal(manifest)
	if err != nil {
		return nil, err
	}

	return &File{
		Name: ManifestFilename,
		Data: data,
	}, nil
}
