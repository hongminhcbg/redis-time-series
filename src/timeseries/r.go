package timeseries

import (
	"time"

	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
)

const days_30 = 30 * 24 * time.Hour

func Init() error {
	// Connect to localhost with no password
	client := redistimeseries.NewClient("localhost:6379", "nohelp", nil)
	keyname := "ytest"
	_, haveit := client.Info(keyname)
	if haveit != nil {
		client.CreateKeyWithOptions(keyname, redistimeseries.DefaultCreateOptions)
		client.CreateKeyWithOptions(keyname+"_avg", redistimeseries.DefaultCreateOptions)
		client.CreateRule(keyname, redistimeseries.AvgAggregation, 60, keyname+"_avg")
	}
	// Add sample with timestamp from server time and value 100
	// TS.ADD mytest * 100
	_, err := client.AddAutoTs(keyname, 100)
	return err
}

type TimeSeries interface {
	AddAutoTs(key string, value float64) error
}

type _ts struct {
	cli *redistimeseries.Client
}

func New() TimeSeries {
	client := redistimeseries.NewClient("localhost:6379", "nohelp", nil)
	return &_ts{cli: client}
}

func (t *_ts) AddAutoTs(key string, value float64) error {
	_, err := t.cli.AddWithOptions(key, time.Now().UnixMilli(), value, redistimeseries.CreateOptions{
		RetentionMSecs: days_30,
	})
	return err
}
