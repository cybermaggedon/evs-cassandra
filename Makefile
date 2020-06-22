
# Create version tag from git tag
VERSION=$(shell git describe | sed 's/^v//')
REPO=cybermaggedon/evs-cassandra
DOCKER=docker
GO=GOPATH=$$(pwd)/go go

all: evs-cassandra build

SOURCE=evs-cassandra.go model.go load.go schema.go config.go

evs-cassandra: ${SOURCE} go.mod go.sum
	${GO} build -o $@ ${SOURCE}

build: evs-cassandra
	${DOCKER} build -t ${REPO}:${VERSION} -f Dockerfile .
	${DOCKER} tag ${REPO}:${VERSION} ${REPO}:latest

push:
	${DOCKER} push ${REPO}:${VERSION}
	${DOCKER} push ${REPO}:latest

