package transformer

import (
	"bytes"
	"strings"
)

type Transformer interface {
	Transform() (*Result, error)
}

type Result struct {
	Buffer   *bytes.Buffer
	Filename string
	FileExt  string
}

func Factory(file string) Transformer {
	if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
		return NewHTTPTransformer(file)
	}
	return NewLocalTransformer(file)
}
