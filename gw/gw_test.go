package gw

import (
	"cti/ds"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DataSourceApiGwTestSuite struct {
	dsApiServer *httptest.Server
	gwApiServer *httptest.Server
	apiApiGw    *DataSourceApiGw
	symbol      string
	suite.Suite
}

func (suite *DataSourceApiGwTestSuite) SetupTest() {
	datasource, err := ds.NewBinanceDataSource()
	suite.Nil(err)

	suite.dsApiServer = httptest.NewServer(ds.NewDataSourceApiServer(datasource, ":8080").Router)

	apiClient, err := ds.NewDefaultDataSourceApiClient(suite.dsApiServer.URL)
	suite.Nil(err)
	suite.symbol = "BTCUSD"

	priceDataSources := []ds.PriceDataSourceApi{apiClient}
	averageDataSources := []ds.AverageDataSourceApi{apiClient}

	suite.apiApiGw = NewDataSourceApiGw(priceDataSources, averageDataSources, suite.symbol, ":8080")
}

func (suite *DataSourceApiGwTestSuite) TestNoDataSource() {
	var priceDataSources []ds.PriceDataSourceApi
	var averageDataSources []ds.AverageDataSourceApi

	apiApiGw := NewDataSourceApiGw(priceDataSources, averageDataSources, suite.symbol, ":8080")

	server := httptest.NewServer(apiApiGw.router)
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/api/v1/price?ts=1569484800")
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(400, resp.StatusCode)

	resp, err = server.Client().Get(server.URL + "/api/v1/average?from=1569484800&until=1569492000&ts=1667457091&granularity=1h")
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(400, resp.StatusCode)
}

func (suite *DataSourceApiGwTestSuite) TestRoutePattern() {
	server := httptest.NewServer(suite.apiApiGw.router)
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/api/v1/price?ts=1569484800")
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(200, resp.StatusCode)

	resp, err = server.Client().Get(server.URL + "/api/v1/average?from=1569484800&until=1569492000&ts=1667457091&granularity=1h")
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(200, resp.StatusCode)
}

func (suite *DataSourceApiGwTestSuite) TestPriceHandler() {
	req, err := http.NewRequest("GET", "/?ts=1569484800", nil)
	suite.Nil(err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.apiApiGw.price)
	handler.ServeHTTP(rr, req)
	suite.Equal(200, rr.Code)
}

func (suite *DataSourceApiGwTestSuite) TestPriceHandlerQsMissingAndUnsupported() {
	handler := http.HandlerFunc(suite.apiApiGw.price)

	var testUrls = []string{
		"/",
		"/?ts=x",
	}

	for _, url := range testUrls {
		req, err := http.NewRequest("GET", url, nil)
		suite.Nil(err)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		suite.Equal(400, rr.Code, "%s, status code: %d", url, rr.Code)
	}

}

func (suite *DataSourceApiGwTestSuite) TestAverageHandler() {
	handler := http.HandlerFunc(suite.apiApiGw.average)

	req, err := http.NewRequest("GET", "/average?from=1569484800&until=1569492000&ts=1667457091&granularity=1h", nil)
	suite.Nil(err)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Equal(200, rr.Code)

	// request without granularity
	req, err = http.NewRequest("GET", "/average?from=1569484800&until=1569492000&ts=1667457091", nil)
	suite.Nil(err)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Equal(200, rr.Code)
}

func (suite *DataSourceApiGwTestSuite) TestAverageHandlerWithNotSupportedGranularity() {
	handler := http.HandlerFunc(suite.apiApiGw.average)

	req, err := http.NewRequest("GET", "/average?from=1569484800&until=1569492000&ts=1667457091&granularity=1x", nil)
	suite.Nil(err)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Equal(400, rr.Code)
}

func (suite *DataSourceApiGwTestSuite) TestAverageHandlerQsMissingAndUnsupported() {
	handler := http.HandlerFunc(suite.apiApiGw.average)

	var testUrls = []string{
		"/",
		"/?from=1569484800",
		"/?until=1569484800",
		"/?from=x&until=1569484800",
		"/?from=x",
		"/?from=1569484800&until=x",
		"/?from=1569484800",
		"/?until=1569484800",
		"/?granularity=1569484800",
	}

	for _, url := range testUrls {
		req, err := http.NewRequest("GET", url, nil)
		suite.Nil(err)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		suite.Equal(400, rr.Code, "%s, status code: %d", url, rr.Code)
	}

}

func (suite *DataSourceApiGwTestSuite) TearDownTestSuite() {
	suite.dsApiServer.Close()
	suite.gwApiServer.Close()
}

func TestDataSourceApiTestSuite(t *testing.T) {
	suite.Run(t, new(DataSourceApiGwTestSuite))
}

func TestErrorPayload(t *testing.T) {
	err := NewErrorPayload(fmt.Errorf("testing"))
	assert.Equal(t, err.Error(), err.Msg)
}
