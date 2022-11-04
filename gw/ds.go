package gw

import "cti/ds"

type DataSourceApiClient interface {
	ds.DataSourceApiClient
	Id() string
}

type DefaultDataSourceApiClient struct {
	ds.DataSourceApiClient
	id string
}

func NewDefaultDataSourceApiClient(id string, client ds.DataSourceApiClient) *DefaultDataSourceApiClient {
	return &DefaultDataSourceApiClient{
		DataSourceApiClient: client,
		id:                  id,
	}
}

func (client DefaultDataSourceApiClient) Id() string {
	return client.id
}
