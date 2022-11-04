package gw

import (
	"cti/ds"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDefaultDataSourceApiClient(t *testing.T) {
	datasource, err := ds.NewDefaultDataSourceApiClient("https://127.0.0.1")
	assert.Nil(t, err)

	id := "id1"
	client := NewDefaultDataSourceApiClient(id, datasource)
	assert.Equal(t, id, client.Id())
}
