package ds

import (
	"cti/db"
	"fmt"
	"github.com/stretchr/testify/suite"
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

type InfluxDbDataSourceTestSuite struct {
	datasource *InfluxDbDataSource
	dataPoints []struct {
		t     time.Time
		value float64
	}
	symbol string
	suite.Suite
}

func (suite *InfluxDbDataSourceTestSuite) SetupTest() {
	datasource, err := NewInfluxDbDataSource(serverUrl, org, bucket, token)
	suite.Nil(err)

	suite.symbol = "BTCUSD"
	suite.datasource = datasource
	suite.dataPoints = []struct {
		t     time.Time
		value float64
	}{
		{t: time.Now().Add(-time.Minute).Truncate(time.Minute), value: 1},
		{t: time.Now().Truncate(time.Minute), value: 2},
	}

	dbWriter := db.NewInfluxDbWriter(serverUrl, org, bucket, token)
	for _, dp := range suite.dataPoints {
		err = dbWriter.WritePrice(suite.symbol, dp.value, dp.t)
		suite.Nil(err)
	}
}

func (suite *InfluxDbDataSourceTestSuite) TestPrice() {
	price, err := suite.datasource.Price(suite.symbol, suite.dataPoints[0].t)
	suite.Nil(err)
	suite.Equal(suite.dataPoints[0].value, price)
}

func (suite *InfluxDbDataSourceTestSuite) TestPriceWithInvalidParams() {
	params := []struct {
		symbol string
		t      time.Time
	}{
		{
			symbol: "",
			t:      time.Time{},
		},
		{
			symbol: suite.symbol,
			t:      time.Time{},
		},
		{
			symbol: suite.symbol,
			t:      time.Now().Add(-time.Second),
		},
		{
			symbol: suite.symbol,
			t:      suite.dataPoints[1].t.Add(time.Hour),
		},
	}

	for _, param := range params {
		_, err := suite.datasource.Price(param.symbol, param.t)
		suite.NotNil(err)
	}
}

func (suite *InfluxDbDataSourceTestSuite) TestAverage() {
	from := suite.dataPoints[0].t
	until := suite.dataPoints[1].t
	fmt.Println("TestAverage", from, until)
	average, actualFrom, actualUntil, err := suite.datasource.Average(suite.symbol, from, until, Granularity1m)
	suite.Nil(err)
	suite.NotEqual((suite.dataPoints[0].value+suite.dataPoints[1].value)/2, average)
	suite.Equal(from.UTC(), actualFrom)
	suite.Equal(until.UTC(), actualUntil)
}

func (suite *InfluxDbDataSourceTestSuite) TestAverageWithInvalidParams() {
	params := []struct {
		symbol      string
		from        time.Time
		until       time.Time
		granularity Granularity
	}{
		{
			symbol:      "",
			from:        time.Time{},
			until:       time.Time{},
			granularity: Granularity1s,
		},
		{
			symbol:      "",
			from:        time.Time{},
			until:       time.Time{},
			granularity: Granularity1m,
		},
		{
			symbol:      suite.symbol,
			from:        suite.dataPoints[0].t,
			until:       time.Time{},
			granularity: Granularity1m,
		},
		{
			symbol:      suite.symbol,
			from:        time.Time{},
			until:       suite.dataPoints[1].t,
			granularity: Granularity1m,
		},
		{
			symbol:      suite.symbol,
			from:        time.Time{},
			until:       time.Time{},
			granularity: Granularity(""),
		},
	}

	for _, param := range params {
		_, _, _, err := suite.datasource.Average(suite.symbol, param.from, param.until, param.granularity)
		suite.NotNil(err)
	}

}

func TestInfluxDbDataSourceTestSuite(t *testing.T) {
	suite.Run(t, new(InfluxDbDataSourceTestSuite))
}
