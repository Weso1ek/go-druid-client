package httpcontroller

import (
	"encoding/json"
	"fmt"
	"github.com/grafadruid/go-druid"
	druidDs "github.com/grafadruid/go-druid/builder/datasource"
	druidGranularity "github.com/grafadruid/go-druid/builder/granularity"
	"github.com/grafadruid/go-druid/builder/query"
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

	d.PrepareDruidQuery(params)

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

func (d DruidController) PrepareDruidQuery(params appContext.InputParams) []byte {
	datasource := os.Getenv("DRUID_DS")

	table := druidDs.NewTable().SetName(datasource)
	granulation := druidGranularity.NewSimple().SetGranularity(params.Granulation)

	ts := druidQuery.NewTimeseries().SetDataSource(table)

	q := `{
		  "queryType": "groupBy",
		  "dataSource": {
			"type": "table",
			"name": "double-sketch"
		  },
		  "granularity": "ALL",
		  "dimensions": [
			{
			  "type": "default",
			  "dimension": "uniqueId"
			}
		  ],
		  "aggregations": [
			{
			  "type": "quantilesDoublesSketch",
			  "name": "a1:agg",
			  "fieldName": "latencySketch",
			  "k": 128
			}
		  ],
		  "postAggregations": [
			{
			  "type": "quantilesDoublesSketchToQuantile",
			  "name": "tp90",
			  "fraction": 0.9,
			  "field": {
				"type": "fieldAccess",
				"name": "tp90",
				"fieldName": "a1:agg"
			  }
			}
		  ],
		  "intervals": {
			"type": "intervals",
			"intervals": [
			  "-146136543-09-08T08:23:32.096Z/146140482-04-24T15:36:27.903Z"
			]
		  }
		}`

	return []byte(q)
}

func (d DruidController) DruidRequest(query []byte) {
	host := os.Getenv("DRUID_HOST")

	client, err := druid.NewClient("http://" + host + "/druid/v2/")

	if err != nil {
		fmt.Println("Druid error ", err)
	}

	fmt.Println(client)

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
