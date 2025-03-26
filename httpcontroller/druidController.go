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
	host := os.Getenv("DRUID_HOST")
	datasource := os.Getenv("DRUID_DS")

	client, err := druid.NewClient("http://" + host + "/druid/v2/")

	if err != nil {
		fmt.Println("Druid error ", err)
	}

	fmt.Println(client)

	table := druidDs.NewTable().SetName(datasource)
	granulation := druidGranularity.NewSimple().SetGranularity(params.Granulation)

	ads := druidAggregation.NewHLLSketchMerge().SetName("uu").SetFieldName("unique")
	a := []druidBuilder.Aggregator{ads}

	filterSite := druidFilter.NewSelector().SetDimension("site").SetValue(strconv.Itoa(params.Site))

	filterBotDimension := druidFilter.NewSelector().SetDimension("bot").SetValue("isrobot")
	filterBot := druidFilter.NewNot().SetField(filterBotDimension)

	filter := druidFilter.NewAnd().

	ts := druidQuery.NewTimeseries().
		SetDataSource(table).
		SetGranularity(granulation).
		SetAggregations(a)

	fmt.Println(ts)

	//q, err := d.Query().Load([]byte(query))
	//if err != nil {
	//	log.Fatalf("converting string to query object failed, %s", err)
	//}
	//
	//resp, err := d.Query().Execute(q, &results)

	client.Close()

}

func (d DruidController) PrepareParams(r *http.Request) appContext.InputParams {
	dateStart, _ := strconv.Atoi(r.URL.Query().Get("dateStart"))
	dateEnd, _ := strconv.Atoi(r.URL.Query().Get("dateEnd"))
	pm, _ := strconv.Atoi(r.URL.Query().Get("pm"))
	pmCategory, _ := strconv.Atoi(r.URL.Query().Get("pmCategory"))
	site, _ := strconv.Atoi(r.URL.Query().Get("site"))

	return appContext.InputParams{
		PmCategory:  pmCategory,
		Pm:          pm,
		Site:        site,
		DateStart:   int32(dateStart),
		DateEnd:     int32(dateEnd),
		Report:      r.URL.Query().Get("report"),
		Granulation: druidGranularity.Day,
		Channel:     r.URL.Query()["channel[]"],
	}
}
