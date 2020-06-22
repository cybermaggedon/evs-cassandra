package main

import (
	"github.com/cybermaggedon/evs-golang-api"
	"os"
	"strings"
)

type CassandraConfig struct {
	*evs.Config
	hosts []string
}

func NewCassandraConfig() *CassandraConfig {

	c := &CassandraConfig{
		Config:  evs.NewConfig("evs-cassandra", "ioc"),
		hosts: []string{"localhost"},
	}

	if val, ok := os.LookupEnv("CASSANDRA_CLUSTER"); ok {
		c.Hosts(strings.Split(val, ","))
	}

	return c

}

func (cc *CassandraConfig) Hosts(hosts []string) {
	cc.hosts = hosts
}
