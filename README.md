aguin (adélie penguin)
=====

For collecting daily/weekly stats or data of different kind. They should be simple data structure of integer, float or boolean. Array of integer/float can be accepted.

The data can be sent signed/encrypted using HMAC/AES or just not if your channel is already secure.


You need to have mongodb installed, http://www.mongodb.org/, for storing the data.

You need Go to run/compile this program http://golang.org/

You can get started by running
	
	make
	make gorun
	
where "make" include

	make godeps
	make goinstall
	make gotestdeps
	make gotest

You can build running
	
	make gobuild

If you want to run "go command" then just export GOPATH to current working directory, such as

	export GOPATH=$(pwd) && go run aguin.go -pid aguin.pid


Init account and add testapp, note down the output (api_key, api_secret, aes_key). You will get the same result if you run for same data again.

	export GOPATH=$(pwd) && go run scripts/init.go -e your.email@fakedomain.tld -n "Your name" -a testapp

Your can send some data by using one of the client lib or send from command line using python client
	
	python clients/python/scripts/cli.py -h

But make sure install dependencies
	
	cd clients/python/
	python setup.py develop