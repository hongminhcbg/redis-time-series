package rules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
)

type VelocityOuput struct {
	UserId  string
	DataKey string
	Ts      string
	Data    float64
}

type RuleMsg struct {
	Input         []byte
	VelocityOuput []*VelocityOuput
	UserId        string

	jsonContainer *gabs.Container
}

func (r *RuleMsg) Parse() error {
	jsonParsed, err := gabs.ParseJSON(r.Input)
	if err != nil {
		return err
	}

	r.jsonContainer = jsonParsed
	return nil
}

func (r *RuleMsg) Extract(key, path string) {
	if strings.HasPrefix(path, "$.") {
		path = path[2:]
	}

	value, ok := r.jsonContainer.Path(path).Data().(float64)
	if !ok {
		return
	}

	timestamp := "*"
	timestampPath := ""
	timestampArgs := strings.Split(path, ".")
	if len(timestampArgs) > 0 {
		timestampArgs[len(timestampArgs)-1] = "timestamp"
		timestampPath = strings.Join(timestampArgs, ".")
	}

	if r.jsonContainer.Path(timestampPath).Exists() {
		if val, ok := r.jsonContainer.Path(timestampPath).Data().(float64); ok {
			timestamp = fmt.Sprintf("%d", int64(val))
		}

		if val, ok := r.jsonContainer.Path(timestampPath).Data().(string); ok {
			if tsInt64, err := strconv.ParseInt(val, 10, 64); err == nil {
				timestamp = fmt.Sprint(tsInt64)
			}
		}
	}

	r.VelocityOuput = append(
		r.VelocityOuput,
		&VelocityOuput{
			UserId:  r.UserId,
			DataKey: key,
			Ts:      timestamp,
			Data:    float64(value),
		},
	)
	return
}

func (r *RuleMsg) GetTransType() string {
	val, ok := r.jsonContainer.Path("trans_type").Data().(string)
	if !ok {
		return ""
	}

	return val
}
