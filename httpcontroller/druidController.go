package httpcontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DruidController struct{}

func NewDruidController() *DruidController {
	return &DruidController{}
}

func (d DruidController) StatDau(w http.ResponseWriter, r *http.Request) {
	report := r.URL.Query()

	fmt.Println(report)

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
