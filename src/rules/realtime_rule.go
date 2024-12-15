package rules

import (
	"fmt"
	"log"
	"time"

	"github.com/hongminhcbg/velocity-rule/src/timeseries"
)

type VelocityDB struct {
	ts      timeseries.TimeSeries
	datakey string
	from    int64
	to      int64
}

func (v *VelocityDB) WithIn(stime int64) *VelocityDB {
	now := time.Now().UnixMilli()
	v.to = now
	v.from = now - stime*1000
	return v
}

func (v *VelocityDB) Sum() float64 {
	log.Printf("v.ts.AvgAggregationSum(%s, %d, %d)\n", v.datakey, v.from, v.to)
	value, err := v.ts.AvgAggregationSum(v.datakey, v.from, v.to)
	if err != nil {
		log.Printf("v.ts.AvgAggregationSum(%s, %d, %d) error: %v\n", v.datakey, v.from, v.to, err)
		return 0
	}

	return value
}

type RealtimeRuleInput struct {
	Ts     timeseries.TimeSeries
	UserId string
	Amount float64
	Msg    string
}

func (r *RealtimeRuleInput) VelocityData(datakey string) *VelocityDB {
	return &VelocityDB{
		ts:      r.Ts,
		datakey: fmt.Sprintf("%s*%s", r.UserId, datakey),
	}
}
