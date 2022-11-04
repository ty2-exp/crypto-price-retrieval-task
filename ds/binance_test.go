package ds

import (
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type BinanceDataSourceTestSuite struct {
	datasource *BinanceDataSource
	symbol     string
	suite.Suite
}

func (suite *BinanceDataSourceTestSuite) SetupTest() {
	datasource, err := NewBinanceDataSource()
	suite.symbol = "BTCUSD"
	suite.Nil(err)
	suite.datasource = datasource
}

func (suite *BinanceDataSourceTestSuite) TestPrice() {
	price, err := suite.datasource.Price(suite.symbol, time.Now().Add(-time.Second).Truncate(time.Second))
	suite.Nil(err)
	suite.NotEqual(0, price)
}

func (suite *BinanceDataSourceTestSuite) TestPriceWithInvalidParams() {
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
			symbol: "",
			t:      time.Now().Add(time.Hour),
		},
		{
			symbol: suite.symbol,
			t:      time.Now().Add(time.Hour),
		},
	}

	for _, param := range params {
		_, err := suite.datasource.Price(param.symbol, param.t)
		suite.NotNil(err)
	}
}

func (suite *BinanceDataSourceTestSuite) TestAverage() {
	from := time.Now().Add(-time.Minute).Truncate(time.Second)
	until := time.Now().Add(-time.Second * 10).Truncate(time.Second)
	average, actualFrom, actualUntil, err := suite.datasource.Average(suite.symbol, from, until, Granularity1s)
	suite.Nil(err)
	suite.NotEqual(0, average)
	suite.Equal(from, actualFrom)
	suite.Equal(until, actualUntil)
}

func (suite *BinanceDataSourceTestSuite) TestAverageWithInvalidParams() {
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
			granularity: Granularity(""),
		},
		{
			symbol:      "",
			from:        time.Time{},
			until:       time.Now().Add(time.Hour).Truncate(time.Second),
			granularity: Granularity1s,
		},
		{
			symbol:      "",
			from:        time.Now().Add(time.Hour).Truncate(time.Second),
			until:       time.Time{},
			granularity: Granularity1s,
		},
		{
			symbol:      "",
			from:        time.Now().Add(time.Hour).Truncate(time.Second),
			until:       time.Time{},
			granularity: Granularity1s,
		},
	}

	for _, param := range params {
		_, _, _, err := suite.datasource.Average(suite.symbol, param.from, param.until, param.granularity)
		suite.NotNil(err)
	}

}

func (suite *BinanceDataSourceTestSuite) TestApiKlines() {
	from := time.Now().Add(-time.Minute).Truncate(time.Second)
	until := time.Now().Add(-time.Second * 10).Truncate(time.Second)
	_, err := suite.datasource.api.Klines(suite.symbol, BinanceApiInterval1s, from.UnixMilli(), until.UnixMilli(), 1)
	suite.Nil(err)
}

func (suite *BinanceDataSourceTestSuite) TestApiKlinesWithInvalidParams() {
	from := time.Now().Add(-time.Minute).Truncate(time.Second)
	until := time.Now().Add(-time.Second * 10).Truncate(time.Second)
	_, err := suite.datasource.api.Klines(suite.symbol, "", from.UnixMilli(), until.UnixMilli(), 1)
	suite.NotNil(err)

	options := []BinanceApiOption{BinanceApiBaseUrlOption("http://fake.x"), BinanceApiHttpClientOption(&http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: 10,
			MaxIdleConns:    100,
		},
		Timeout: time.Millisecond,
	})}
	dataSouce, err := NewBinanceDataSource(options...)
	suite.Nil(err)

	_, err = dataSouce.api.Klines(suite.symbol, BinanceApiInterval1d, from.UnixMilli(), until.UnixMilli(), 1)
	suite.NotNil(err)
}

func (suite *BinanceDataSourceTestSuite) TestWithOption() {
	options := []BinanceApiOption{BinanceApiBaseUrlOption("http://127.0.0.1"), BinanceApiHttpClientOption(http.DefaultClient)}
	_, err := NewBinanceDataSource(options...)
	suite.Nil(err)

	options2 := []BinanceApiOption{BinanceApiBaseUrlOption("http://user:abc{DEf1=ghi@127.0.0.1")}
	_, err = NewBinanceDataSource(options2...)
	suite.NotNil(err)

	options3 := []BinanceApiOption{BinanceApiBaseUrlOption("httpsx://127.0.0.1")}
	_, err = NewBinanceDataSource(options3...)
	suite.NotNil(err)

}

func TestBinanceDataSourceTestSuite(t *testing.T) {
	suite.Run(t, new(BinanceDataSourceTestSuite))
}
