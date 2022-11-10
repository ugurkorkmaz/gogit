package gogit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	git := &git{}
	gogitNew := New()

	assert.Equal(t, git, gogitNew)
}
