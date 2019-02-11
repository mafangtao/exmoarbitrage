package route

import (
	"github.com/tusupov/exmoarbitrage/config"
	"github.com/tusupov/exmoarbitrage/route/controller"
	"github.com/tusupov/exmoarbitrage/service"

	"github.com/gorilla/mux"
)

func Init(cfg *config.Config, service service.Servicer) (router *mux.Router, err error) {

	web, err := controller.NewWeb(cfg, service)
	if err != nil {
		return
	}

	router = mux.NewRouter()
	router.HandleFunc("/", web.Index)
	router.HandleFunc("/arbitrage", web.Arbitrage)
	router.HandleFunc("/currency", web.Currency)

	return

}
