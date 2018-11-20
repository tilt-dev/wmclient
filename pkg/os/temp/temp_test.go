package temp

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTempDir(t *testing.T) {
	dir, err := NewDir("prefix")
	assert.NoError(t, err)
	assert.Contains(t, filepath.Base(dir.Path()), "prefix")
}

func TestTempDirError(t *testing.T) {
	_, err := NewDir("in/valid")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "mkdir")
	}
}
