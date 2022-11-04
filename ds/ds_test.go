package ds

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGranularity(t *testing.T) {
	g := Granularity1s
	assert.Equal(t, true, g.IsValid())

	invalidG := Granularity("")
	assert.Equal(t, false, invalidG.IsValid())
}
