<?php

class Aguin {
	const VERSION = 'AguinPHP/0.1';	
	const ENCRYPT_METHOD = 'AES-256-CFB';
	
	private $apiKey;
	private $apiSecret;
	private $aesKey;
	private $url;
	private $https;
	
	public function __construct($apiKey, $apiSecret, $aesKey, $url) {
		$this->apiKey = $apiKey;
		$this->apiSecret = $apiSecret;
		$this->aesKey = $aesKey;
		$this->url = rtrim($url, '/');
		$parsedUrl = parse_url($url);
		$this->https = $parsedUrl['scheme'] == 'https';
	}
	
	public function get($entity, array $criteria = NULL) {
		if ($criteria == NULL) {
			$criteria = array();
		}
		$criteria['entity'] = $entity;
		$message = $this->encrypt($criteria);
		return $this->request('GET', $this->url, array('message' => $message));
	}
	
	public function post($entity, array $data) {
		$message = $this->encrypt(array('entity' => $entity, 'data' => $data));
		return $this->request('POST', $this->url, array('message' => $message));
	}
	
	public function status() {
		return $this->request('GET', $this->url . '/status');	
	}
	
	public function decrypt($base64String) {
		$parts = explode('.', $base64String);
		if (count($parts) != 2) {
			throw new Exception('Invalid data');
		}
		$expectedHmac = base64_decode($this->urlsafeDecode($parts[0]));
		$message = base64_decode($this->urlsafeDecode($parts[1]));
		
		if ($expectedHmac != hash_hmac('sha256', $message , $this->apiSecret, TRUE)) {
			throw new Exception('Invalid data');
		}
		
		return json_decode($this->aesDecrypt($message), TRUE);
	}
	
	public function encrypt(array $data) {
		$message = $this->aesEncrypt(json_encode($data));
		$expectedHmac = $this->urlsafeEncode(base64_encode(hash_hmac('sha256', $message , $this->apiSecret, TRUE)));
		return $expectedHmac . '.' . $this->urlsafeEncode(base64_encode($message));
	}
	
	private function request($method, $url, array $data = NULL) {
		$headers = array('X-AGUIN-API-KEY: ' . $this->apiKey);
		
		$ch = curl_init();
		
		curl_setopt($ch, CURLOPT_USERAGENT, self::VERSION);
		curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, TRUE);
		
		if ($method == 'POST') {
			curl_setopt($ch, CURLOPT_POST, 1);
			curl_setopt($ch, CURLOPT_POSTFIELDS, http_build_query($data));
		} else {
			if ($data != NULL) {
				$url = $url . '?'. http_build_query($data);
			}
		}
		
		curl_setopt($ch, CURLOPT_URL, $url);
		
		if ($this->https) {
			curl_setopt($ch, CURLOPT_PROTOCOLS, CURLPROTO_HTTPS);
		}
		
		if (!($result = curl_exec($ch))) {
			throw new Exception(curl_error($ch));
		}
		
		$jresult = json_decode($result, TRUE);
		
		if ($jresult['encrypted']) {
			$jresult['result'] = $this->decrypt($jresult['result']);
		}
		
		return $jresult;
	}
	
	private function aesDecrypt($message) {
		$iv = substr($message, 0, openssl_cipher_iv_length(self::ENCRYPT_METHOD));
		$data = substr($message, openssl_cipher_iv_length(self::ENCRYPT_METHOD));
		return openssl_decrypt($data, self::ENCRYPT_METHOD, $this->aesKey, OPENSSL_RAW_DATA, $iv);
	}
	
	private function aesEncrypt($message) {
		$iv = openssl_random_pseudo_bytes(openssl_cipher_iv_length(self::ENCRYPT_METHOD));
		$encrypted = openssl_encrypt($message, self::ENCRYPT_METHOD, $this->aesKey, OPENSSL_RAW_DATA, $iv);
		return $iv . $encrypted;
	}
	
	private function urlsafeDecode($base64String) {
		return str_replace(array('-', '_'), array('+','/'), $base64String);
	}
	
	private function urlsafeEncode($base64String) {
		return str_replace(array('/','+'), array('_', '-'), $base64String);
	}
}


//$api = new Aguin('545e0716f2fea0c7a9c46c74', '545e0716f2fea0c7a9c46c74fec46c71', '545e0716f2fea0c7a9c46c74fec46c71', 'http://127.0.0.1:8080/');
//print $api->decrypt('KxdeCmADcSszX7X6PniP6n12jf6pv7VhPCfevZGR48k=.ny3maaaunErzWNutsQQK67aWXgYwHDtcSMhgun5SaDSl-GKWdCll');
//print_r($api->status());
//print_r($api->post('something', array('php' => 1, 'a' => 10, 'b' => 0.2)));
//print_r($api->get('something'));


