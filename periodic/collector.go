package periodic

import (
	"cti/db"
	"cti/ds"
	"fmt"
	"log"
	"sync"
	"time"
)

type Collector struct {
	datasource ds.DataSourceApiClient
	dbWriter   db.Writer
	symbol     string
	stopCh     chan struct{}
	startOnce  *sync.Once
	stopOnce   *sync.Once
	interval   time.Duration
}

func NewCollector(datasource ds.DataSourceApiClient, dbWriter db.Writer, symbol string) *Collector {
	collector := &Collector{
		datasource: datasource,
		dbWriter:   dbWriter,
		symbol:     symbol,
		interval:   time.Minute,
	}

	collector.reset()

	return collector
}

func (collector *Collector) reset() {
	collector.stopCh = make(chan struct{}, 1)
	collector.startOnce = &sync.Once{}
	collector.stopOnce = &sync.Once{}
}

func (collector *Collector) Start() {
	collector.startOnce.Do(collector.start)
}

func (collector *Collector) start() {
	err := collector.collect()
	if err != nil {
		log.Println(err)
	}

	ticker := time.NewTicker(collector.interval)
	for {
		select {
		case <-ticker.C:
			err := collector.collect()
			if err != nil {
				log.Println(err)
				continue
			}
		case <-collector.stopCh:
			break
		}
	}
}

func (collector *Collector) collect() error {
	ts := time.Now().Add(-collector.interval).Truncate(collector.interval)

	price, err := collector.datasource.Price(collector.symbol, ts)
	if err != nil {
		return fmt.Errorf("price request fail: %w", err)
	}

	err = collector.dbWriter.WritePrice(collector.symbol, price.Price, ts)
	if err != nil {
		return fmt.Errorf("db write fail: %w", err)
	}

	return nil
}

func (collector *Collector) Stop() {
	collector.stopOnce.Do(collector.stop)
}

func (collector *Collector) stop() {
	collector.stopCh <- struct{}{}
	close(collector.stopCh)
	collector.reset()
}
