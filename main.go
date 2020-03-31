package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/w564791/prometheus-cert-exporter/cert"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	//"k8s.io/client_go/rest"
	//"math/rand"
	"net/http"
	"os"
)

var (
	Paths = kingpin.Flag("path", "cert path,provide file/dir").Required().Strings()
)

type ClusterManager struct {
	Zone        string
	CERTmanager *prometheus.Desc
}

//func (c *ClusterManager) ReallyExpensiveAssessmentOfTheSystemState() (
//	oomCountByHost []map[string]interface{},
//	//ramUsageByHost map[string]float64,
//) {
//
//	ToomCountByHost := map[string]interface{}{
//		"value": float64(rand.Int31n(1000)),
//		"label": []string{"2019","2018"},
//	}
//	//ramUsageByHost = map[string]float64{
//	//	"foo.example.org": rand.Float64() * 100,
//	//	"bar.example.org": rand.Float64() * 100,
//	//}
//	oomCountByHost= append(oomCountByHost, ToomCountByHost)
//	return
//}

func (c *ClusterManager) ReallyExpensiveAssessmentOfTheSystemState(path string) (
	lables []string, value float64,
	//ramUsageByHost map[string]float64,
) {

	//var ss=make(map[string]string)
	lables, value = cert.ParsePem(path)
	//log.Error("ssss",lables)
	//oomCountByHost=append(oomCountByHost, map[string]interface{}{
	//	"value":value,
	//	"labels":lables,
	//
	//})

	//ramUsageByHost = map[string]float64{
	//	"foo.example.org": rand.Float64() * 100,
	//	"bar.example.org": rand.Float64() * 100,
	//}
	return
}
func (c *ClusterManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.CERTmanager

}

func (c *ClusterManager) Collect(ch chan<- prometheus.Metric) {
	//var oomCountByHost []string
	var labels []string
	var value float64
	for _, path := range *Paths {
		if cert.IsFile(path) {
			//log.Error("err labels value",*labels,*value)
			labels, value = c.ReallyExpensiveAssessmentOfTheSystemState(path)
			if labels == nil {
				return
			}
			ch <- prometheus.MustNewConstMetric(
				c.CERTmanager,
				prometheus.CounterValue,
				value,
				labels...,
			)

		} else if cert.IsDir(path) {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				log.Error(err.Error())
			}
			for _, file := range files {
				labels, value = c.ReallyExpensiveAssessmentOfTheSystemState(fmt.Sprintf("%s/%s", path, file.Name()))
				if labels == nil {
					return
				}
				ch <- prometheus.MustNewConstMetric(
					c.CERTmanager,
					prometheus.CounterValue,
					value,
					labels...,
				)

			}
		}

		//log.Error("6666",*value,*labels)

		//label:=strings.Join(host,",")

	}

	//var ss=make(map[string]string)
	//var from,  =
	//for _,p :=range {
	//	from,after,date:=cert.ParsePem(p,ss)
	//
	//}
	//log.Error("oomCountByHost:",oomCountByHost)

}

func NewClusterManager(zone string) *ClusterManager {
	return &ClusterManager{
		Zone: zone,
		CERTmanager: prometheus.NewDesc(
			"cert_exp_date",
			"cert_exp_date",
			[]string{"after", "from", "name"},
			prometheus.Labels{},
		),
	}
}

func main() {
	kingpin.Parse()
	workerDB := NewClusterManager("db")

	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(workerDB)
	//reg.MustRegister(workerCA)
	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		reg,
	}

	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
	log.Infoln("Start server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Errorf("Error occur when start server %v", err)
		os.Exit(1)
	}

}
