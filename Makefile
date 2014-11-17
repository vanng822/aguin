
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
	go install aguin/utils
	go install aguin/config
	go install aguin/crypto
	go install aguin/validator
	go install aguin/model
	go install aguin/api
	

gotest:
	go test aguin/crypto
	go test aguin/validator
