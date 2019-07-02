package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"math/rand"
	"net/http"
	"os"
)


var (
	cpuTEMP = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature",
		Help: "Current temperature of the CPU",
	})
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors",
		},
		[]string{"device"},
	)
	flag1=1.0
	flag2=1.0
)

func init() {
	//Metrics have to be registered to be exposed;
	prometheus.MustRegister(cpuTEMP)
	prometheus.MustRegister(hdFailures)
}

type ClusterManager struct {
	Zone         string
	OOMCountDesc *prometheus.Desc
	RAMUsageDesc *prometheus.Desc
	JiaLaoShi *prometheus.Desc
	DongJie *prometheus.Desc
	SongLaoShi *prometheus.Desc
}

func NewClusterManager(zone string) *ClusterManager {
	return &ClusterManager{
		Zone: zone,
		OOMCountDesc: prometheus.NewDesc(
			"clustermanager_oom_crashes_total",
			"Number of OOM crashes",
			[]string{"host"},
			prometheus.Labels{"zone": zone},
		),
		RAMUsageDesc: prometheus.NewDesc(
			"clustermanager_ram_usage_bytes",
			"RAM usage as reported to the cluster manager",
			[]string{"host"},
			prometheus.Labels{"zone": zone},
		),
		JiaLaoShi: prometheus.NewDesc(
			"JiaLaoShiUpUp",
			"MANY MANY JiaLaoShi",
			[]string{"JiaLaoShiYa"},
			prometheus.Labels{"zone":zone},
			),
			DongJie:prometheus.NewDesc(
				"WuDiDeDongJie",
				"MANY MANY Dong Jie",
				[]string{"DongJieDongJie"},
				prometheus.Labels{"zone":zone},
				),
			SongLaoShi:prometheus.NewDesc(
				"SongLaoShiXuan",
				"MANY MANY SongLaoShi",
				[]string{"SongLaoShiYa"},
				prometheus.Labels{"zone":zone},
				),
	}
}

func (c *ClusterManager) ReallyExpensiveAssessmentOfTheSystemState() (
	oomCountByHost map[string]int, ramUsageByHost map[string]float64,jiajia map[string]int,songsong map[string]float64,dongdong map[string]float64,
) {
	flag2=flag2*1.1
	flag1=flag2*rand.Float64()
	oomCountByHost = map[string]int{
		"foo.example.org": int(rand.Int31n(1000)),
		"bar.example.org": int(rand.Int31n(1000)),
	}
	ramUsageByHost = map[string]float64{
		"foo.example.org": rand.Float64() * 100,
		"bar.example.org": rand.Float64() * 100,
	}
	jiajia= map[string]int{
		"JiaYuHang01": int(rand.Int31n(10000)),
		"JiaYuHang02": int(rand.Int31n(8888)),
	}
	songsong=map[string]float64{
		"SLY01": flag1*rand.Float64()*100,
		"SlY02": flag1*rand.Float64()*100,
	}
	songsong=map[string]float64{
		"SLY01": flag1*rand.Float64()*100,
		"SlY02": flag1*rand.Float64()*100,
	}
	dongdong= map[string]float64{
		"DongJie01": flag2*rand.Float64()*100,
		"DongJie02": flag2*rand.Float64()*100,
	}
	return
}

func (c *ClusterManager) Describe(ch chan<- *prometheus.Desc) {//一定要注册！
	ch <- c.OOMCountDesc
	ch <- c.RAMUsageDesc
	ch <- c.JiaLaoShi
	ch <- c.SongLaoShi
	ch <- c.DongJie
}

func (c *ClusterManager) Collect(ch chan<- prometheus.Metric) {
	oomCountByHost, ramUsageByHost,jiaLaoShi,songLaoShi,dongJie := c.ReallyExpensiveAssessmentOfTheSystemState()
	for host, oomCount := range oomCountByHost {
		ch <- prometheus.MustNewConstMetric(
			c.OOMCountDesc,
			prometheus.CounterValue,
			float64(oomCount),
			host,
		)
	}
	for host, ramUsage := range ramUsageByHost {
		ch <- prometheus.MustNewConstMetric(
			c.RAMUsageDesc,
			prometheus.GaugeValue,
			ramUsage,
			host,
		)
	}
	for host,jiajiajia:=range jiaLaoShi{
		ch <-prometheus.MustNewConstMetric(
			c.JiaLaoShi,
			prometheus.CounterValue,
			float64(jiajiajia),
			host,
			)
	}
	for host,songsongsong:=range songLaoShi{
		ch <- prometheus.MustNewConstMetric(
			c.SongLaoShi,
			prometheus.CounterValue,
			float64(songsongsong),
			host,
			)
	}
	for host,dongdong:=range dongJie{
		ch <- prometheus.MustNewConstMetric(
			c.DongJie,
			prometheus.CounterValue,
			float64(dongdong),
			host,
		)
	}
}

func main() {

	//The Handler function provides a default handler to expose metrics
	//via an HTTP server. "/metrics" is the usual endpoint for that.
	//http.Handle("/", promhttp.Handler())
	//log.Fatal(http.ListenAndServe(":8888", nil))

	workerDB := NewClusterManager("db")
	workerCA := NewClusterManager("ca")

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(workerDB)
	reg.MustRegister(workerCA)

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		reg,
	}
	h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError,
	})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
	log.Infoln("Start server at 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		//log.Errorf("Error occur when start server %v", err)
		os.Exit(1)
	}
}
