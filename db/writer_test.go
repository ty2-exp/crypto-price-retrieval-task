package db

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	serverUrl = os.Getenv("TEST_INFLUX_SERVER_URL")
	org       = os.Getenv("TEST_INFLUX_ORG")
	token     = os.Getenv("TEST_INFLUX_TOKEN")
	bucket    = os.Getenv("TEST_INFLUX_BUCKET")
)

func TestInfluxDbWriterWritePrice(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	dbWriter := NewInfluxDbWriter(serverUrl, org, bucket, token)
	err := dbWriter.WritePrice("BTCUSD", 1, now)
	assert.Nil(t, err)
}

func TestInfluxDbWriterWritePriceWithError(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	dbWriter := NewInfluxDbWriter("", org, bucket, token)
	err := dbWriter.WritePrice("BTCUSD", 1, now)
	assert.NotNil(t, err)
}
