package pkpass

import (
	"io/ioutil"
	"os"
	"os/exec"
)

// NewOpenSSL returns a signer backed by openssl.
func NewOpenSSL(wwdrcert, passcert, passkey, password string) Signer {
	return &openssl{
		WWDRCertificate: wwdrcert,
		PassCertificate: passcert,
		PassKey:         passkey,
		Password:        password,
	}
}

type openssl struct {
	WWDRCertificate string
	PassCertificate string
	PassKey         string
	Password        string
}

func (o *openssl) Sign(data []byte) (*File, error) {
	inFile, err := tempFile(data)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	outFile, err := tempFile([]byte{})
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	cmd := exec.Command(
		"openssl",
		"smime", "-binary",
		"-sign",
		"-certfile", o.WWDRCertificate,
		"-signer", o.PassCertificate,
		"-inkey", o.PassKey,
		"-in", inFile.Name(),
		"-out", outFile.Name(),
		"-outform", "DER",
		"-passin", "pass:"+o.Password,
	)

	if err = cmd.Run(); err != nil {
		return nil, err
	}

	signatureData, err := ioutil.ReadFile(outFile.Name())
	if err != nil {
		return nil, err
	}

	return &File{
		Name: SignatureFilename,
		Data: signatureData,
	}, nil
}

func tempFile(data []byte) (*os.File, error) {
	file, err := ioutil.TempFile("", "okpock")
	if err != nil {
		return nil, err
	}

	_, err = file.Write(data)
	if err != nil {
		return nil, err
	}

	return file, nil
}
