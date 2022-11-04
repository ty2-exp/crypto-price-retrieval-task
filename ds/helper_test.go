package ds

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlParseWithJoin(t *testing.T) {
	u, err := UrlParseWithJoin("http://127.0.0.1", "hello")
	assert.Nil(t, err)
	assert.Equal(t, "http://127.0.0.1/hello", u.String())
}

func TestUrlParseWithJoinInvalid(t *testing.T) {
	_, err := UrlParseWithJoin("http://user:abc{DEf1=ghi@127.0.0.1", "hello")
	assert.NotNil(t, err)
}

func TestEscapeDoubleQuote(t *testing.T) {
	s := EscapeDoubleQuote("\"hello\"")
	assert.Equal(t, `\"hello\"`, s)
}
