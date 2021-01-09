package transformer

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
)

var _ Transformer = (*LocalTransformer)(nil)

type LocalTransformer struct {
	filePath string
}

func NewLocalTransformer(filePath string) *LocalTransformer {
	return &LocalTransformer{filePath: filePath}
}

func (t *LocalTransformer) Transform() (*Result, error) {
	c, err := ioutil.ReadFile(t.filePath)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(c)
	fullName := filepath.Base(t.filePath)
	ext := filepath.Ext(fullName)
	filename := fullName[:len(fullName)-len(ext)]
	return &Result{
		Buffer:   buffer,
		Filename: filename,
		FileExt:  ext,
	}, nil
}
