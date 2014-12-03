require "openssl"
require "digest"
require "base64"
require 'json'
require 'net/https'

VERSION = "Aguin_ruby/0.1"

class Aguin
  @@AES_BLOCK_SIZE = 16
  def initialize(api_key, api_secret, aes_key, url)
    @api_key = api_key
    @api_secret = api_secret
    @aes_key = aes_key
    @url = url.chomp('/')
  end

  def get(entity, criteria = nil)
    if criteria == nil
      criteria = {}
    end
    criteria["entity"] = entity
    uri = URI.parse(@url)
    uri.query = URI.encode_www_form({:message => encrypt(criteria)})
    req = Net::HTTP::Get.new(uri)
    return request(req, uri)
  end

  def post(entity, data)
    uri = URI.parse(@url)
    req = Net::HTTP::Post.new(uri)
    req.set_form_data('message' => encrypt({:entity => entity, :data => data}))
    return request(req, uri)
  end

  def status()
    uri = URI.parse(@url + '/status')
    req = Net::HTTP::Get.new(uri)
    return request(req, uri)
  end
  
  def encrypt(data)
    encrypted = aes_encrypt(JSON.generate(data))
    expected_mac = Base64.urlsafe_encode64(OpenSSL::HMAC::digest(OpenSSL::Digest::SHA256.new, @api_secret, encrypted))
    message = Base64.urlsafe_encode64(encrypted)
    return expected_mac + '.' + message
  end

  def decrypt(base64String)
    parts = base64String.split('.')
    if parts.size != 2
    raise ArgumentError, "Invalid data"
    end
    expected_mac = Base64.urlsafe_decode64(parts[0])
    message = Base64.urlsafe_decode64(parts[1])
    message_mac = OpenSSL::HMAC::digest(OpenSSL::Digest::SHA256.new, @api_secret, message)
    if message_mac != expected_mac
    raise ArgumentError, "Invalid data"
    end
    return JSON.parse(aes_decrypt(message), {:quirks_mode => true})
  end

  private

  def request(req, uri)
    req.initialize_http_header({"User-Agent" => VERSION, "X-AGUIN-API-KEY" => @api_key})
    http = Net::HTTP.new(uri.host, uri.port)
    if uri.instance_of? URI::HTTPS
    http.use_ssl = true
    end
    res = http.request(req)
    bjson = JSON.parse(res.body)

    if bjson["encrypted"]
    bjson["result"] = decrypt(bjson["result"])
    end
    return bjson
  end

  def aes_decrypt(message)
    decipher = OpenSSL::Cipher::AES256.new(:CFB)
    decipher.decrypt
    decipher.key = @aes_key
    decipher.iv = message.slice(0, @@AES_BLOCK_SIZE)
    return decipher.update(message.slice(@@AES_BLOCK_SIZE, message.size)) + decipher.final
  end

  def aes_encrypt(message)
    cipher = OpenSSL::Cipher::AES256.new(:CFB)
    cipher.encrypt
    cipher.key = @aes_key
    iv = cipher.random_iv
    cipher.iv = iv
    return iv + cipher.update(message) + cipher.final
  end
end

t = Aguin::new('545e0716f2fea0c7a9c46c74', '545e0716f2fea0c7a9c46c74fec46c71', '545e0716f2fea0c7a9c46c74fec46c71', 'http://127.0.0.1:8080/')
#puts t.decrypt('KxdeCmADcSszX7X6PniP6n12jf6pv7VhPCfevZGR48k=.ny3maaaunErzWNutsQQK67aWXgYwHDtcSMhgun5SaDSl-GKWdCll')
puts t.status()
#puts t.post('something', {:a => 1, :b => 2, :c => 0.1})
puts t.get('something')
