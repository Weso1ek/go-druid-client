package httpcontroller

import (
	"encoding/json"
	"fmt"
	appContext "go-druid-client/context"
	"net/http"
	"strconv"
)

type DruidController struct{}

func NewDruidController() *DruidController {
	return &DruidController{}
}

func (d DruidController) StatDau(w http.ResponseWriter, r *http.Request) {
	params := d.PrepareParams(r)

	fmt.Println(params)

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

func (d DruidController) PrepareParams(r *http.Request) appContext.InputParams {
	dateStart, _ := strconv.Atoi(r.URL.Query().Get("dateStart"))
	dateEnd, _ := strconv.Atoi(r.URL.Query().Get("dateEnd"))

	return appContext.InputParams{
		PmCategory: 100,
		DateStart:  int32(dateStart),
		DateEnd:    int32(dateEnd),
	}
}
