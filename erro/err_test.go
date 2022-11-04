package erro

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError(t *testing.T) {
	e := NewError("TEST", "test", nil)
	assert.Equal(t, e.Error(), e.Text)
	assert.Equal(t, map[string]interface{}(nil), e.Attr)

	e1, ok := e.WithAttrs(map[string]any{"field": "testField"}).(*Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, e.Text, e1.Text)
	assert.Equal(t, e.Code, e1.Code)
	assert.Equal(t, e.Error(), e1.Error())
	assert.Equal(t, map[string]any{"field": "testField"}, e1.Attr)

	assert.Equal(t, map[string]interface{}(nil), e.Attr)

}
