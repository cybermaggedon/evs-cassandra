//
// ElasticSearch loader
//

package main

import (
	"github.com/gocql/gocql"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

type Loader struct {
	CassandraConfig

	session *gocql.Session
	insert  *gocql.Query

	event_latency prometheus.Summary
	loaded        prometheus.Summary
}

func (l *Loader) InitMetrics() {

	//configuration specific to prometheus stats
	l.event_latency = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "cassandra_event_latency",
			Help: "Latency from probe to store",
		})
	prometheus.MustRegister(l.event_latency)

	l.loaded = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "cassandra_loaded",
			Help: "Numer of Cassandra rows loaded",
		})
	prometheus.MustRegister(l.loaded)

}

func (l *Loader) InitClient() {

	for {

		cluster := gocql.NewCluster()
		cluster.Hosts = l.hosts
		cluster.Consistency = gocql.Quorum

		var err error
		l.session, err = cluster.CreateSession()
		if err != nil {
			log.Printf("Cassandra connection: %v", err)
			time.Sleep(1 * time.Second * 5)
			continue
		}

		break

	}

}

func (l *Loader) Init() error {

	l.InitMetrics()
	l.InitClient()

	err := l.InitSchema()
	if err != nil {
		return err
	}

	return nil

}

func (l *Loader) Load(ob *Event) error {

	qry := l.insert.Bind(ob.Id, ob.Time, ob.Action, ob.Device, ob.Network,
		ob.SrcIp, ob.DestIp, ob.SrcPort, ob.DestPort, ob.Protocol,
		ob.DnsType, ob.DnsQuery, ob.DnsAnswer, ob.HttpRequest,
		ob.HttpResponse, ob.Url, ob.Header, ob.Indicator)

	if err := qry.Exec(); err != nil {
		return err
	}

	ts := time.Now()
	go l.recordLatency(ts, ob)

	return nil

}

func (l *Loader) recordLatency(ts time.Time, ob *Event) {
	latency := ts.Sub(ob.Time)
	l.event_latency.Observe(float64(latency))
}

func (cc CassandraConfig) NewLoader() (*Loader, error) {
	l := &Loader{
		CassandraConfig: cc,
	}
	err := l.Init()
	if err != nil {
		return nil, err
	}
	return l, nil
}
