package pkpass

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/fullsailor/pkcs7"
	"golang.org/x/crypto/pkcs12"
)

// ErrInvalidRootCert returned when pem block is nil.
var ErrInvalidRootCert = errors.New("pkpass: invalid root certificate")

// Signer holds method working with certificates.
type Signer interface {
	Sign(data []byte) (*File, error)
}

// NewSigner returns a new instance of pkcs7 signer.
func NewSigner(root, cert []byte, pass string) (Signer, error) {
	block, _ := pem.Decode(root)
	if block == nil {
		return nil, ErrInvalidRootCert
	}

	rootCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey, signingCert, err := pkcs12.Decode(cert, pass)
	if err != nil {
		return nil, err
	}

	return &signer{
		RootCert:    rootCert,
		SigningCert: signingCert,
		PrivateKey:  privateKey,
	}, nil
}

type signer struct {
	RootCert    *x509.Certificate
	SigningCert *x509.Certificate
	PrivateKey  interface{}
}

func (s *signer) Sign(data []byte) (*File, error) {
	signedData, err := pkcs7.NewSignedData(data)
	if err != nil {
		return nil, err
	}

	signedData.AddCertificate(s.RootCert)

	err = signedData.AddSigner(s.SigningCert, s.PrivateKey, pkcs7.SignerInfoConfig{})
	if err != nil {
		return nil, err
	}

	signedData.Detach()

	signatureData, err := signedData.Finish()
	if err != nil {
		return nil, err
	}

	return &File{
		Name: SignatureFilename,
		Data: signatureData,
	}, nil
}
