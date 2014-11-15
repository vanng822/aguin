import unittest

from aguin import Aguin


"""
def main():
    client = Aguin('545e0716f2fea0c7a9c46c74', '545e0716f2fea0c7a9c46c74fec46c71', '545e0716f2fea0c7a9c46c74fec46c71', 'http://127.0.0.1:8080/')
    #print client.decrypt('KxdeCmADcSszX7X6PniP6n12jf6pv7VhPCfevZGR48k=.ny3maaaunErzWNutsQQK67aWXgYwHDtcSMhgun5SaDSl-GKWdCll')
    #print client.post('something', {'testing': 1, 'hej': 2})
    print client.get('something')
    print client.status()
    
if __name__ == '__main__':
    main()
"""

class TestAguin(unittest.TestCase):
    
    def setUp(self):
        self.client =  Aguin('545e0716f2fea0c7a9c46c74', '545e0716f2fea0c7a9c46c74fec46c71', '545e0716f2fea0c7a9c46c74fec46c71', 'http://127.0.0.1:8080/')
        
    def test_encrypt(self):
        expected = {u'entity': u'something2'}
        result = self.client.encrypt(expected)
        self.assertDictEqual(expected, self.client.decrypt(result))
    
    def test_decrypt(self):
        expected = {u'entity': u'something2'}
        result = self.client.decrypt('KxdeCmADcSszX7X6PniP6n12jf6pv7VhPCfevZGR48k=.ny3maaaunErzWNutsQQK67aWXgYwHDtcSMhgun5SaDSl-GKWdCll')
        self.assertDictEqual(expected, result)