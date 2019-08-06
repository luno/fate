package fate

import (
	"github.com/prometheus/client_golang/prometheus"
)

var temptCount = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "fate",
	Name:      "tempt_total",
	Help:      "Total number of times fate was tempted",
})
var temptErrors = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "fate",
	Name:      "tempt_error_total",
	Help:      "Total number of times fate was tempted and errored",
})

func init() {
	prometheus.MustRegister(temptCount)
	prometheus.MustRegister(temptErrors)
}
