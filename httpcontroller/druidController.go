package httpcontroller

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v11"
	appcontext "go-druid-client/context"
	"net/http"
)

type DruidController struct{}

func NewDruidController() *DruidController {
	return &DruidController{}
}

func (d DruidController) StatDau(w http.ResponseWriter, r *http.Request) {
	report := r.URL.Query().Get("report")

	config := appcontext.Config{}
	err := env.Parse(&config)

	fmt.Println(report)
	fmt.Println("asdasdsa22222222")
	fmt.Println(config.DruidDatasource)
	fmt.Println(config.DruidHost)
	fmt.Println(config.DruidPort)
	fmt.Println("asdasdsa222222222")

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
