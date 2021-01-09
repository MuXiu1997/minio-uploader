package transformer

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/flytam/filenamify"
)

var _ Transformer = (*HTTPTransformer)(nil)

type HTTPTransformer struct {
	rawURL string
}

func NewHTTPTransformer(rawURL string) *HTTPTransformer {
	return &HTTPTransformer{rawURL: rawURL}
}

func (t *HTTPTransformer) Transform() (*Result, error) {
	content, err := downloadFile(t.rawURL)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(content)
	ext := getExt(content)
	fullName := filepath.Base(t.rawURL)
	filename := fullName
	if strings.HasSuffix(filename, ext) {
		filename = fullName[:len(fullName)-len(ext)]
	}
	filename, _ = filenamify.Filenamify(filename, filenamify.Options{})

	return &Result{
		Buffer:   buffer,
		Filename: filename,
		FileExt:  ext,
	}, nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func getExt(content []byte) string {
	var ext string
	contentType := http.DetectContentType(content)
	extensions, err := mime.ExtensionsByType(contentType)
	if err == nil && 0 < len(extensions) {
		ext = extensions[0]
	}
	return ext
}
