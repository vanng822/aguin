aguin (adélie penguin)
=====

For collecting daily/weekly stats of different kind. They should be simple data structure of integer, float or boolean. Array of integer/float can be accepted.

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
