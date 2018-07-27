package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDetectExt(t *testing.T) {

	tests := []struct {
		filename string
		data     []byte
		ext      string
	}{
		{"my.json", []byte("{}"), ".json"},
		{"/bla/bla/bla/my.json", []byte("{}"), ".json"},
		{"my.yaml", []byte("---"), ".yaml"},
		{"my.yml", []byte("---"), ".yml"},
		{"-", []byte(`{\n"step":{}}`), ".json"},
		{"-", []byte(`{"step":{}}`), ".json"},
		{"noext", []byte(`{}`), ".json"},
		{"-", []byte(`---\n`), ".yaml"},
		{"-", []byte(yamlExample), ".yaml"},
		{"my.sh", []byte("#anything here"), ".sh"},
	}

	for _, tc := range tests {
		ext := detectExt(tc.filename, tc.data)
		assert.Equal(t, tc.ext, ext)
	}

}

const yamlExample = `node:
  id: status
`
