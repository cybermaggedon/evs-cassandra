package main

import (
	evs "github.com/cybermaggedon/evs-golang-api"
	pb "github.com/cybermaggedon/evs-golang-api/protos"
	"log"
)

const ()

type Cassandra struct {
	*CassandraConfig

	// Embed event analytic framework
	*evs.EventSubscriber
	evs.Interruptible

	loader *Loader
}

// Initialisation
func NewCassandra(cc *CassandraConfig) *Cassandra {

	c := &Cassandra{
		CassandraConfig: cc,
	}

	var err error
	c.EventSubscriber, err = evs.NewEventSubscriber(c, c)
	if err != nil {
		log.Fatal(err)
	}

	c.RegisterStop(c)

	c.loader, err = c.NewLoader()
	if err != nil {
		log.Fatal(err)
	}

	return c
}

// Event handler for new events.
func (c *Cassandra) Event(ev *pb.Event, p map[string]string) error {

	obs := Convert(ev)

	err := c.loader.Load(obs)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	cc := NewCassandraConfig()
	c := NewCassandra(cc)
	log.Print("Initialisation complete")
	c.Run()
	log.Print("Shutdown.")

}
