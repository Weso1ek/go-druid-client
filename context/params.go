package context

import (
	druidGranularity "github.com/grafadruid/go-druid/builder/granularity"
)

type InputParams struct {
	Report      string
	Pm          int
	PmCategory  int
	Site        int
	Category    []string
	Source      []string
	Channel     []string
	DateStart   int32
	DateEnd     int32
	Granulation druidGranularity.Simple
}
