package periodic

import (
	"cti/db"
	"cti/ds"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
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

type CollectorTestSuite struct {
	collector *Collector
	apiServer *httptest.Server
	apiClient ds.DataSourceApiClient
	symbol    string
	suite.Suite
}

func (suite *CollectorTestSuite) SetupTest() {
	datasource, err := ds.NewBinanceDataSource()
	suite.Nil(err)

	suite.apiServer = httptest.NewServer(ds.NewDataSourceApiServer(datasource, ":8080").Router)

	apiClient, err := ds.NewDefaultDataSourceApiClient(suite.apiServer.URL)
	suite.Nil(err)
	suite.apiClient = apiClient
	suite.symbol = "BTCUSD"
}

func (suite *CollectorTestSuite) TestCollector() {
	apiClient, err := ds.NewDefaultDataSourceApiClient(suite.apiServer.URL)
	suite.Nil(err)

	dbWriter := db.NewInfluxDbWriter(serverUrl, org, bucket, token)
	collector := NewCollector(apiClient, dbWriter, suite.symbol)
	collector.interval = time.Second
	go collector.Start()
	time.Sleep(time.Second * 5)
	collector.Stop()
}

func (suite *CollectorTestSuite) TestCollectorCollect() {
	dbWriter := db.NewInfluxDbWriter(serverUrl, org, bucket, token)
	collector := NewCollector(suite.apiClient, dbWriter, suite.symbol)
	err := collector.collect()
	suite.Nil(err)
}

func (suite *CollectorTestSuite) TearDownTestSuite() {
	suite.apiServer.Close()
}

func TestBinanceDataSourceTestSuite(t *testing.T) {
	suite.Run(t, new(CollectorTestSuite))
}
