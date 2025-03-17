package httpcontroller

import (
	"encoding/json"
	"fmt"
	appContext "go-druid-client/context"
	"net/http"
)

type DruidController struct{}

func NewDruidController() *DruidController {
	return &DruidController{}
}

func (d DruidController) StatDau(w http.ResponseWriter, r *http.Request) {
	params := d.PrepareParams(r.URL.Query())

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

func (d DruidController) PrepareParams(query map[string][]string) appContext.InputParams {

	fmt.Println(query)

	return appContext.InputParams{
		PmCategory: 100,
	}
}
