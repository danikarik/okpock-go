package pkpass

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
)

// ErrEmptyFolder returned when there is no files to be zipped.
var ErrEmptyFolder = errors.New("pkpass: no files given to be zipped")

// NewFile returns a new instance of `File`.
func NewFile(name string, data []byte) File {
	return File{name, data}
}

// File holds name and content of pkpass' file.
type File struct {
	Name string
	Data []byte
}

// Zip archives files into zip.
func Zip(files ...File) ([]byte, error) {
	if len(files) == 0 {
		return nil, ErrEmptyFolder
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, file := range files {
		zipFile, err := zipWriter.Create(file.Name)
		if err != nil {
			return nil, err
		}

		_, err = zipFile.Write(file.Data)
		if err != nil {
			return nil, err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Unzip opens zip package.
func Unzip(data []byte) ([]File, error) {
	if len(data) == 0 {
		return nil, ErrEmptyFolder
	}

	num := int64(len(data))
	buf := bytes.NewReader(data)

	zipReader, err := zip.NewReader(buf, num)
	if err != nil {
		return nil, err
	}

	files := make([]File, len(zipReader.File))
	for i, zipFile := range zipReader.File {
		r, err := zipFile.Open()
		if err != nil {
			return nil, err
		}
		defer r.Close()

		data, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		files[i] = File{
			Name: zipFile.Name,
			Data: data,
		}
	}

	return files, nil
}
