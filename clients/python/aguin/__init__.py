
import json
import base64
import hmac, hashlib
import M2Crypto
import requests

VERSION='aguin_py/0.1'

class Aguin(object):
    def __init__(self, api_key, api_secret, aes_key, url):
        self.api_key = api_key
        self.api_secret = api_secret
        self.aes_key = aes_key
        self.url = url.rstrip('/')

    def encrypt(self, data):
        encrypted = self.aes_encrypt(json.dumps(data))
        
        expected_mac = base64.urlsafe_b64encode(hmac.new(self.api_secret, encrypted, hashlib.sha256).digest())
        message = base64.urlsafe_b64encode(encrypted)
        return '{expected_mac}.{message}'.format(expected_mac=expected_mac, message=message)
    
    def decrypt(self, base64String):
        parts = base64String.split('.')
        if len(parts) != 2:
            raise Exception('encrypted data was in wrong format')
        
        expected_mac = base64.urlsafe_b64decode(parts[0])
        message = base64.urlsafe_b64decode(parts[1])
        message_mac = hmac.new(self.api_secret, message, hashlib.sha256).digest()
        if message_mac != expected_mac:
            raise Exception('Invalid data')
        
        return json.loads(self.aes_decrypt(message))
    
    def aes_decrypt(self, message):
        if len(message) < M2Crypto.m2.AES_BLOCK_SIZE:
            raise Exception('Message is too short')
        
        return m2_decrypt(message[M2Crypto.m2.AES_BLOCK_SIZE:], self.aes_key, message[:M2Crypto.m2.AES_BLOCK_SIZE])

    def aes_encrypt(self, data):
        iv = M2Crypto.m2.rand_bytes(16)
        message = m2_encrypt(data, self.aes_key, iv)
    
        return iv + message
    
    def post(self, entity, data):
        encrypted = self.encrypt({'entity': entity, 'data': data})
        return self.request(params={'message': encrypted}, method='POST')
    
    def get(self, entity, criteria=None):
        if criteria is None:
            data = {}
        else:
            data = criteria
            
        data['entity'] = entity
        return self.request(params={'message': self.encrypt(data)})
    
    def status(self):
        return self.request(path='/status')
        
    def request(self, path='/',params=None, method='GET'):
        url = '{url}{path}'.format(url=self.url, path=path)
        if params is None:
            params = {}
        
        headers = {'X-AGUIN-API-KEY': self.api_key, 'User-Agent': VERSION}
        if method == 'POST':
            res = requests.post(url, params=params, headers=headers)
        else:
            res = requests.get(url, params=params, headers=headers)
        result = res.json()
        
        if result['encrypted'] == True:
            result['result'] = self.decrypt(str(result['result']))
            
        return result

def m2_encrypt(plaintext, key, iv, key_as_bytes=False, padding=True):
    cipher = M2Crypto.EVP.Cipher(alg="aes_256_cfb", key=key, iv=iv,
                        key_as_bytes=key_as_bytes, padding=padding, op=1)
    return cipher.update(plaintext) + cipher.final()
 
def m2_decrypt(ciphertext, key, iv, key_as_bytes=False, padding=True):
    cipher = M2Crypto.EVP.Cipher(alg="aes_256_cfb", key=key, iv=iv,
                        key_as_bytes=key_as_bytes, padding=padding, op=0)
    return cipher.update(ciphertext) + cipher.final()