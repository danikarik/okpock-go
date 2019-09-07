package pkpass

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/fullsailor/pkcs7"
	"golang.org/x/crypto/pkcs12"
)

// SignatureFilename is an alias for `signature`.
const SignatureFilename = "signature"

// Signer holds method working with certificates.
type Signer interface {
	Sign(data []byte) (*File, error)
}

const wwdr = `
-----BEGIN CERTIFICATE-----
MIIEIjCCAwqgAwIBAgIIAd68xDltoBAwDQYJKoZIhvcNAQEFBQAwYjELMAkGA1UE
BhMCVVMxEzARBgNVBAoTCkFwcGxlIEluYy4xJjAkBgNVBAsTHUFwcGxlIENlcnRp
ZmljYXRpb24gQXV0aG9yaXR5MRYwFAYDVQQDEw1BcHBsZSBSb290IENBMB4XDTEz
MDIwNzIxNDg0N1oXDTIzMDIwNzIxNDg0N1owgZYxCzAJBgNVBAYTAlVTMRMwEQYD
VQQKDApBcHBsZSBJbmMuMSwwKgYDVQQLDCNBcHBsZSBXb3JsZHdpZGUgRGV2ZWxv
cGVyIFJlbGF0aW9uczFEMEIGA1UEAww7QXBwbGUgV29ybGR3aWRlIERldmVsb3Bl
ciBSZWxhdGlvbnMgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQDKOFSmy1aqyCQ5SOmM7uxfuH8mkbw0U3rOfGOA
YXdkXqUHI7Y5/lAtFVZYcC1+xG7BSoU+L/DehBqhV8mvexj/avoVEkkVCBmsqtsq
Mu2WY2hSFT2Miuy/axiV4AOsAX2XBWfODoWVN2rtCbauZ81RZJ/GXNG8V25nNYB2
NqSHgW44j9grFU57Jdhav06DwY3Sk9UacbVgnJ0zTlX5ElgMhrgWDcHld0WNUEi6
Ky3klIXh6MSdxmilsKP8Z35wugJZS3dCkTm59c3hTO/AO0iMpuUhXf1qarunFjVg
0uat80YpyejDi+l5wGphZxWy8P3laLxiX27Pmd3vG2P+kmWrAgMBAAGjgaYwgaMw
HQYDVR0OBBYEFIgnFwmpthhgi+zruvZHWcVSVKO3MA8GA1UdEwEB/wQFMAMBAf8w
HwYDVR0jBBgwFoAUK9BpR5R2Cf70a40uQKb3R01/CF4wLgYDVR0fBCcwJTAjoCGg
H4YdaHR0cDovL2NybC5hcHBsZS5jb20vcm9vdC5jcmwwDgYDVR0PAQH/BAQDAgGG
MBAGCiqGSIb3Y2QGAgEEAgUAMA0GCSqGSIb3DQEBBQUAA4IBAQBPz+9Zviz1smwv
j+4ThzLoBTWobot9yWkMudkXvHcs1Gfi/ZptOllc34MBvbKuKmFysa/Nw0Uwj6OD
Dc4dR7Txk4qjdJukw5hyhzs+r0ULklS5MruQGFNrCk4QttkdUGwhgAqJTleMa1s8
Pab93vcNIx0LSiaHP7qRkkykGRIZbVf1eliHe2iK5IaMSuviSRSqpd1VAKmuu0sw
ruGgsbwpgOYJd+W+NKIByn/c4grmO7i77LpilfMFY0GCzQ87HUyVpNur+cmV6U/k
TecmmYHpvPm0KdIBembhLoz2IYrF+Hjhga6/05Cdqa3zr/04GpZnMBxRpVzscYqC
tGwPDBUf
-----END CERTIFICATE-----`

// NewSigner returns a new instance of pkcs7 signer.
func NewSigner(couponCert string) Signer {
	return &signer{couponCert}
}

type signer struct {
	couponCert string
}

func (s *signer) Sign(data []byte) (*File, error) {
	signedData, err := pkcs7.NewSignedData(data)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(wwdr))
	if block == nil {
		return nil, errors.New("could not load wwdr certificate")
	}

	wwdrCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	signedData.AddCertificate(wwdrCert)

	couponData, err := ioutil.ReadFile(s.couponCert)
	if err != nil {
		return nil, err
	}

	couponPrivKey, couponCert, err := pkcs12.Decode(couponData, "dthcnf07")
	if err != nil {
		return nil, err
	}

	err = signedData.AddSigner(couponCert, couponPrivKey, pkcs7.SignerInfoConfig{})
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
