package httpcontroller

import (
	"encoding/json"
	"fmt"
	"github.com/grafadruid/go-druid"
	druidBuilder "github.com/grafadruid/go-druid/builder"
	druidAggregation "github.com/grafadruid/go-druid/builder/aggregation"
	druidDs "github.com/grafadruid/go-druid/builder/datasource"
	druidFilter "github.com/grafadruid/go-druid/builder/filter"
	druidGranularity "github.com/grafadruid/go-druid/builder/granularity"
	"log"
	"time"

	druidIn "github.com/grafadruid/go-druid/builder/intervals"
	druidQuery "github.com/grafadruid/go-druid/builder/query"
	appContext "go-druid-client/context"
	"net/http"
	"os"
	"strconv"
)

type DruidController struct{}

func NewDruidController() *DruidController {
	return &DruidController{}
}

func (d DruidController) StatDau(w http.ResponseWriter, r *http.Request) {
	params := d.PrepareParams(r)

	d.DruidRequest(params)

	var status = map[string]string{
		"status": "ok",
	}

	statusJson, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(statusJson)
}

func (d DruidController) DruidRequest(params appContext.InputParams) {
	client, err := druid.NewClient("https://" + os.Getenv("DRUID_HOST") + ":" + os.Getenv("DRUID_PORT"))

	if err != nil {
		fmt.Println("Druid error ", err)
	}

	table := druidDs.NewTable().SetName(os.Getenv("DRUID_DS"))
	granulation := druidGranularity.NewSimple().SetGranularity(params.Granulation)

	ads := druidAggregation.NewHLLSketchMerge().SetName("uu").SetFieldName("unique")
	a := []druidBuilder.Aggregator{ads}

	// site condition
	filterSite := druidFilter.NewSelector().SetDimension("site").SetValue(strconv.Itoa(params.Site))

	// bot condition
	filterBotDimension := druidFilter.NewSelector().SetDimension("bot").SetValue("isrobot")
	filterBot := druidFilter.NewNot().SetField(filterBotDimension)

	// pos filter
	filterPosEmpty := druidFilter.NewSelector().SetDimension("pos").SetValue("")
	filterPosT := druidFilter.NewSelector().SetDimension("pos").SetValue("T")
	filterPosF := druidFilter.NewSelector().SetDimension("pos").SetValue("F")

	filterPos := druidFilter.NewOr().SetFields([]druidBuilder.Filter{filterPosEmpty, filterPosT, filterPosF})

	// filter channel
	filterChannel1 := druidFilter.NewSelector().SetDimension("channel").SetValue("1")
	filterChannel2 := druidFilter.NewSelector().SetDimension("channel").SetValue("2")
	filterChannel3 := druidFilter.NewSelector().SetDimension("channel").SetValue("3")

	filterChannel := druidFilter.NewOr().SetFields([]druidBuilder.Filter{filterChannel1, filterChannel2, filterChannel3})

	filter := druidFilter.NewAnd().SetFields([]druidBuilder.Filter{filterSite, filterBot, filterPos, filterChannel})

	// intervals
	dateStart := params.DateStartTime
	dateEnd := params.DateEndTime

	interval := druidIn.NewInterval().SetIntervalWithString(dateStart.Format(time.RFC3339Nano), dateEnd.Format(time.RFC3339Nano))
	intervals := druidIn.NewIntervals().SetIntervals([]*druidIn.Interval{interval})

	ts := druidQuery.NewTimeseries().
		SetDataSource(table).
		SetGranularity(granulation).
		SetAggregations(a).
		SetIntervals(intervals).
		SetFilter(filter)

	var results interface{}

	_, err = client.Query().Execute(ts, &results)

	if err != nil {
		log.Fatalf("Execute failed, %s", err)
	}

	client.Close()
}

func (d DruidController) PrepareParams(r *http.Request) appContext.InputParams {
	dateStart, _ := strconv.Atoi(r.URL.Query().Get("dateStart"))
	dateEnd, _ := strconv.Atoi(r.URL.Query().Get("dateEnd"))
	pm, _ := strconv.Atoi(r.URL.Query().Get("pm"))
	pmCategory, _ := strconv.Atoi(r.URL.Query().Get("pmCategory"))
	site, _ := strconv.Atoi(r.URL.Query().Get("site"))

	return appContext.InputParams{
		PmCategory:    pmCategory,
		Pm:            pm,
		Site:          site,
		DateStart:     int64(dateStart),
		DateStartTime: time.Unix(int64(dateStart), 0),
		DateEnd:       int64(dateEnd),
		DateEndTime:   time.Unix(int64(dateEnd), 0),
		Report:        r.URL.Query().Get("report"),
		Granulation:   druidGranularity.Day,
		Channel:       r.URL.Query()["channel[]"],
	}
}
