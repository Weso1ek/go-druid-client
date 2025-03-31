package context

import (
	druidGranularity "github.com/grafadruid/go-druid/builder/granularity"
	"time"
)

type InputParams struct {
	Report        string
	Pm            int
	PmCategory    int
	Site          int
	Category      []string
	Source        []string
	Channel       []string
	DateStart     int64
	DateStartTime time.Time
	DateEnd       int64
	DateEndTime   time.Time
	Granulation   druidGranularity.Simple
}
