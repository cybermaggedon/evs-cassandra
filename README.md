# `evs-cassandra`

Eventstream analytic for Cyberprobe event streams.  Subscribes to Pulsar
for Cyberprobe events and produces events which are stored in a table on
Cassandra.

## Getting Started

The target deployment product is a container engine.  The analytic expects
a Pulsar service to be running, along with a Cassandra service.

```
  docker run -d \
      -e PULSAR_BROKER=pulsar://<PULSAR-HOST>:6650 \
      -e CASSANDRA_CLUSTER=<CASSANDRA-HOST1>,... \
      -p 8088:8088 \
      docker.io/cybermaggedon/evs-cassandra:<VERSION>
```
      
### Prerequisites

You need to have a container deployment system e.g. Podman, Docker, Moby.

You need to have a Cassandra service running.  To run a single-node,
non-production Cassandra service, try...

```
  docker run -d -p 9042:9042 cassandra:3.11.6
```

You also need a Pulsar exchange, being fed by events from Cyberprobe.

### Installing

The easiest way is to use the containers we publish to Docker hub.
See https://hub.docker.com/r/cybermaggedon/evs-cassandra

```
  docker pull docker.io/cybermaggedon/evs-cassandra:<VERSION>
```

If you want to build this yourself, you can just clone the Github repo,
and type `make`.

## Deployment configuration

The following environment variables are used to configure:

| Variable | Purpose | Default |
|----------|---------|---------|
| `INPUT` | Specifies the Pulsar topic to subscribe to.  This is just the topic part of the URL e.g. `cyberprobe`. | `ioc` |
| `METRICS_PORT` | Specifies the port number to serve Prometheus metrics on.  If not set, metrics will not be served. The container has a default setting of 8088. | `8088` |
| `CASSANDRA_CLUSTER` | Specifies a set of contact points used to discover the Cassandra cluster topology.  Should be a comma-separated list.  Use a single hostname if the cluster is a single node. | `localhost` |


