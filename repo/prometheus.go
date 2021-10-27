package repo

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type prometheusRepository struct {

}

var (
	prometheusImpl *prometheusRepository
)

func GetPrometheusRepo() *prometheusRepository {
	if prometheusImpl == nil {
		prometheusImpl = NewPrometheusRepo()
	}
	return prometheusImpl
}

func NewPrometheusRepo() *prometheusRepository {
	var prom = new(prometheusRepository)
	return prom
}

func (repo *prometheusRepository)init()*prometheusRepository {
	return repo
}

func (repo *prometheusRepository)GetHttpHandler() http.Handler {
	 return promhttp.Handler()
}