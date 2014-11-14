Aguin = require('../');

t = new Aguin('545e0716f2fea0c7a9c46c74', '545e0716f2fea0c7a9c46c74fec46c71', '545e0716f2fea0c7a9c46c74fec46c71', 'http://127.0.0.1:8080/')


t.post("something", {"test": 1, "t": 0.1}, function(err, result) {
	console.log(err, result);
});

t.get({"entity":"something"}, function(err, result) {
	console.log(err, result);
});
