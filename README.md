aguin (adélie penguin)
=====

For collecting daily/weekly stats of different kind. They should be simple data structure of date (2006-01-02 03:04:01), integer, float or boolean. Array of integer/float can be accepted.

The data can be sent signed/encrypted using HMAC/AES or just not if your channel is already secure.


You need to have mongodb installed, http://www.mongodb.org/, for storing the data.

You need Go to run/compile this program http://golang.org/

You need to install packages

	go get github.com/go-martini/martini
	
	go get github.com/martini-contrib/render
	
	go get gopkg.in/mgo.v2
