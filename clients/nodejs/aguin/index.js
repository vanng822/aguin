
var crypto = require('crypto');
var http = require('http');
var https = require('https');
var url = require('url');


var Aguin = function(apiKey, apiSecret, aesSecret, url) {
	this.apiKey = apiKey;
	this.apiSecret = apiSecret;
	this.aesSecret = aesSecret;
	this.url = url;
}


Aguin.prototype = {
	encrypt : function(data) {
		var iv, cipher, encoded, hmac, expectedMac;
		iv = crypto.randomBytes(16);
		cipher = crypto.createCipheriv('aes-256-cfb', new Buffer(this.aesSecret), iv);
		cipher.setAutoPadding(true);
		encoded = Buffer.concat([iv, cipher.update(JSON.stringify(data), 'utf8'), cipher.final()]);
		hmac = crypto.createHmac('sha256', this.apiSecret);
		expectedMac = hmac.update(encoded).digest();
		return this.urlencode(expectedMac.toString('base64')) + '.' + this.urlencode(encoded.toString('base64'));
	},
	decrypt : function(encryptedData, callback) {
		var parts, expectedMac, message, messageMAC;
		var hmac, decipher;
		parts = String(encryptedData).split('.');
		if (parts.length != 2) {
			throw Error('Invalid data');
		}
		expectedMac = new Buffer(this.urldecode(parts[0]), 'base64');
		message = new Buffer(this.urldecode(parts[1]), 'base64');
		hmac = crypto.createHmac('sha256', this.apiSecret);
		messageMac = hmac.update(message).digest();
		// Buffer and SlowBuffer!
		if (messageMac.toString() != expectedMac.toString()) {
			throw Error('Invalid data');
		}
		decipher = crypto.createDecipheriv('aes-256-cfb', new Buffer(this.aesSecret), message.slice(0, 16));
		decipher.setAutoPadding(true);
		
		return decipher.update(message.slice(16)) + decipher.final('utf8');
	},
	post: function(entity, data, callback) {
		var options = url.parse(this.url);
		options.method = 'POST';
		options.headers = {'X-AGUIN-API-KEY': this.apiKey};
		options.path = '/?message=' + encodeURIComponent(this.encrypt({entity:entity, data: data}));
		this.request(options, callback);
	},
	get: function(entity, criteria, callback) {
		var options = url.parse(this.url);
		var data;
		options.method = 'GET';
		options.headers = {'X-AGUIN-API-KEY': this.apiKey};
		if (!callback) {
			callback = criteria;
			criteria = null;
		}
		if (criteria) {
			data = criteria;
		} else {
			data = {};
		}
		data.entity = entity;
		options.path = '/?message=' + encodeURIComponent(this.encrypt(data));
		this.request(options, callback);
	},
	request : function(options, callback) {
		var h, self = this;
		if (options.protocol === 'https:') {
			h = https;
		} else {
			h = http;
		}
		console.log(options);
		var req = h.request(options, function(res) {
			var data = '';
			res.on('data', function(chunk) {
				data += chunk;
			});
			res.on('end', function() {
				var result;
				try {
					result = JSON.parse(data);
					if (result.encrypted) {
						result.result = JSON.parse(self.decrypt(result.result));
						callback(null, result);
					} else {
						callback(null, result);
					}
					
				} catch(e) {/* json parse may throw error */
					return callback(e, null);
				}
			});
		}).on('error', function(err) {
			callback(err, null);
		});
		req.end();
	},
	urlencode : function(base64String) {
		return base64String.replace('+', '-').replace('/', '_');
	},
	urldecode: function(base64String) {
		return base64String.replace('-', '+').replace('_', '/');
	}
};

module.exports = Aguin