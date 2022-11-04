package ds

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type DataSourceApiTestSuite struct {
	datasource DataSource
	apiServer  *DataSourceApiServer
	suite.Suite
}

func (suite *DataSourceApiTestSuite) SetupTest() {
	datasource, err := NewBinanceDataSource()
	suite.Nil(err)
	suite.datasource = datasource

	suite.apiServer = NewDataSourceApiServer(suite.datasource, ":8080")
}

func (suite *DataSourceApiTestSuite) TestRoutePattern() {
	server := httptest.NewServer(suite.apiServer.Router)
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/api/v1/price?symbol=BTCUSD&ts=1569484800")
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(200, resp.StatusCode)

	resp, err = server.Client().Get(server.URL + "/api/v1/average?symbol=BTCUSD&from=1569484800&until=1569492000&ts=1667457091&granularity=1h")
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(200, resp.StatusCode)
}

func (suite *DataSourceApiTestSuite) TestPriceHandler() {
	req, err := http.NewRequest("GET", "/?symbol=BTCUSD&ts=1569484800", nil)
	suite.Nil(err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.apiServer.Price)
	handler.ServeHTTP(rr, req)
	suite.Equal(200, rr.Code)
}

func (suite *DataSourceApiTestSuite) TestPriceHandlerQsMissingAndUnsupported() {
	handler := http.HandlerFunc(suite.apiServer.Price)

	var testUrls = []string{
		"/",
		"/?symbol=BTCUSD",
		"/?ts=1569484800",
		"/?ts=x",
		"/?symbol=BTCUSD&ts=x",
	}

	for _, url := range testUrls {
		req, err := http.NewRequest("GET", url, nil)
		suite.Nil(err)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		suite.Equal(400, rr.Code, "%s, status code: %d", url, rr.Code)
	}

}

func (suite *DataSourceApiTestSuite) TestAverageHandler() {
	handler := http.HandlerFunc(suite.apiServer.Average)

	req, err := http.NewRequest("GET", "/average?symbol=BTCUSD&from=1569484800&until=1569492000&ts=1667457091&granularity=1h", nil)
	suite.Nil(err)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Equal(200, rr.Code)

	// request without granularity
	req, err = http.NewRequest("GET", "/average?symbol=BTCUSD&from=1569484800&until=1569492000&ts=1667457091", nil)
	suite.Nil(err)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Equal(200, rr.Code)
}

func (suite *DataSourceApiTestSuite) TestAverageHandlerWithNotSupportedGranularity() {
	handler := http.HandlerFunc(suite.apiServer.Average)

	req, err := http.NewRequest("GET", "/average?symbol=BTCUSD&from=1569484800&until=1569492000&ts=1667457091&granularity=1x", nil)
	suite.Nil(err)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Equal(400, rr.Code)
}

func (suite *DataSourceApiTestSuite) TestAverageHandlerQsMissingAndUnsupported() {
	handler := http.HandlerFunc(suite.apiServer.Average)

	var testUrls = []string{
		"/",
		"/?symbol=BTCUSD",
		"/?symbol=BTCUSD&from=1569484800",
		"/?symbol=BTCUSD&until=1569484800",
		"/?symbol=BTCUSD&from=x&until=1569484800",
		"/?symbol=BTCUSD&from=x",
		"/?symbol=BTCUSD&from=1569484800&until=x",
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

func TestDataSourceApiTestSuite(t *testing.T) {
	suite.Run(t, new(DataSourceApiTestSuite))
}

type DataSourceApiClientTestSuite struct {
	datasource DataSource
	apiServer  *DataSourceApiServer
	apiClient  *DataSourceApiServer
	symbol     string
	suite.Suite
}

func (suite *DataSourceApiClientTestSuite) SetupTest() {
	datasource, err := NewBinanceDataSource()
	suite.Nil(err)
	suite.datasource = datasource
	suite.symbol = "BTCUSD"

	suite.apiServer = NewDataSourceApiServer(suite.datasource, ":8080")
}

func (suite *DataSourceApiClientTestSuite) TestPrice() {
	server := httptest.NewServer(suite.apiServer.Router)
	defer server.Close()
	apiClient, err := NewDefaultDataSourceApiClient(server.URL)
	suite.Nil(err)

	result, err := apiClient.Price(suite.symbol, time.Now().Add(-time.Minute))
	suite.Nil(err)
	suite.NotEqual(0, result.Price)
}

func (suite *DataSourceApiClientTestSuite) TestPriceWithoutSymbol() {
	server := httptest.NewServer(suite.apiServer.Router)
	defer server.Close()
	apiClient, err := NewDefaultDataSourceApiClient(server.URL)
	suite.Nil(err)

	result, err := apiClient.Price("", time.Now().Add(-time.Minute))
	suite.NotNil(err)

	suite.NotEqual(0, result.Price)
}

func (suite *DataSourceApiClientTestSuite) TestAverage() {
	server := httptest.NewServer(suite.apiServer.Router)
	defer server.Close()
	apiClient, err := NewDefaultDataSourceApiClient(server.URL)
	suite.Nil(err)

	result, err := apiClient.Average(suite.symbol, time.Now().Add(-time.Minute*5).Truncate(time.Minute), time.Now().Add(-time.Minute*2), Granularity1m)
	suite.Nil(err)
	suite.NotEqual(0, result.Average)
}

func (suite *DataSourceApiClientTestSuite) TestAverageWithoutSymbol() {
	server := httptest.NewServer(suite.apiServer.Router)
	defer server.Close()
	apiClient, err := NewDefaultDataSourceApiClient(server.URL)
	suite.Nil(err)

	result, err := apiClient.Average("", time.Now().Add(-time.Minute*5).Truncate(time.Minute), time.Now().Add(-time.Minute*2).Truncate(time.Minute), Granularity1m)
	suite.NotNil(err)
	suite.NotEqual(0, result.Average)
}

func (suite *DataSourceApiClientTestSuite) TestAverageWithCustomClient() {
	server := httptest.NewServer(suite.apiServer.Router)
	defer server.Close()
	_, err := NewDefaultDataSourceApiClient(server.URL, []DefaultDataSourceApiClientOption{DefaultDataSourceApiClientHttpClientOption(http.DefaultClient)}...)
	suite.Nil(err)
}

func TestDataSourceApiClientTestSuite(t *testing.T) {
	suite.Run(t, new(DataSourceApiClientTestSuite))
}

func TestErrorPayload(t *testing.T) {
	err := NewErrorPayload(fmt.Errorf("testing"))
	assert.Equal(t, err.Error(), err.Msg)
}
