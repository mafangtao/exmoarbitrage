package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"github.com/tusupov/exmoarbitrage/config"
	"github.com/tusupov/exmoarbitrage/service"
)

type Web struct {
	cfg     *config.Config
	service service.Servicer

	indexTpl     *template.Template
	arbitrageTpl *template.Template
	currencyTpl  *template.Template
}

func NewWeb(cfg *config.Config, service service.Servicer) (web *Web, err error) {

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}

	t, err := template.New("index.html").Funcs(funcMap).ParseFiles(cfg.TemplateDirectory + "view/index.html")
	if err != nil {
		return
	}
	indexTpl := template.Must(t, nil)

	t, err = template.New("arbitrage.html").Funcs(funcMap).ParseFiles(cfg.TemplateDirectory + "view/arbitrage.html")
	if err != nil {
		return
	}
	arbitrageTpl := template.Must(t, nil)

	t, err = template.New("currency.html").Funcs(funcMap).ParseFiles(cfg.TemplateDirectory + "view/currency.html")
	if err != nil {
		return
	}
	currencyTpl := template.Must(t, nil)

	web = &Web{
		cfg:          cfg,
		service:      service,
		indexTpl:     indexTpl,
		arbitrageTpl: arbitrageTpl,
		currencyTpl:  currencyTpl,
	}

	return

}

func (c *Web) Index(w http.ResponseWriter, r *http.Request) {

	arbitrageList, err := c.service.GetArbitrage(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list := make([][]interface{}, 0)
	for _, arbitrage := range arbitrageList {
		list = append(
			list,
			[]interface{}{
				fmt.Sprintf("%.4f", (arbitrage.Profit-1)*100),
				fmt.Sprint(arbitrage.Route),
			},
		)
		if len(list) == 10 {
			break
		}
	}

	err = c.indexTpl.Execute(w, map[string]interface{}{
		"url":  r.URL.String(),
		"list": list,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c *Web) Arbitrage(w http.ResponseWriter, r *http.Request) {

	arbitrageList, err := c.service.GetArbitrage(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list := make([][]interface{}, 0)
	for _, arbitrage := range arbitrageList {
		list = append(
			list,
			[]interface{}{
				fmt.Sprintf("%.4f", (arbitrage.Profit-1)*100),
				fmt.Sprint(arbitrage.Route),
			},
		)
	}

	err = c.arbitrageTpl.Execute(w, map[string]interface{}{
		"url":  r.URL.String(),
		"list": list,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (c *Web) Currency(w http.ResponseWriter, r *http.Request) {

	currencyList, err := c.service.GetCurrencyList(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.currencyTpl.Execute(w, map[string]interface{}{
		"url":  r.URL.String(),
		"list": currencyList,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
