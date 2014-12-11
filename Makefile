export GOPATH := $(shell pwd)

packages=$(shell ls ./src/aguin/ )

all:
	make godeps
	make goinstall
	make gotestdeps
	make gotest

godeps:
	go get github.com/go-martini/martini
	go get github.com/martini-contrib/render
	go get gopkg.in/mgo.v2

gotestdeps:
	go get github.com/stretchr/testify/assert

goinstall:
	$(foreach package, ${packages}, go install aguin/$(package) ;)
	

gotest:
	$(foreach package, ${packages}, go test aguin/$(package) ;)
	
gorun:
	make goinstall
	go run aguin.go -pid aguin.pid
	
	
gobuild:
	make all
	go build aguin.go