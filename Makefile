all:
	make goinstall
	make gotest

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
