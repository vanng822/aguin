import argparse
import json

from aguin import Aguin

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('-k', '--api-key', required=True, type=str, help='Api key')
    parser.add_argument('-s', '--api-secret', required=True, type=str, help='Api secret')
    parser.add_argument('-ak', '--aes-key', required=True, type=str, help='AES key')
    parser.add_argument('-u', '--url', required=True, type=str, help='Url to api')
    parser.add_argument('-a', '--action', required=True, type=str, choices=['GET','POST'], help='Indicate if POST or GET')
    parser.add_argument('-e', '--entity', required=True, type=str, help='Name of the entity')
    parser.add_argument('-d', '--data', type=str, help='Json data to send. If action is post then it is entity data, if get then search criteria')
    parser.add_argument('-p', '--pretty', action='store_true', help='Pretty print of the result')
    
    args = parser.parse_args()
    
    api = Aguin(args.api_key, args.api_secret, args.aes_key, args.url)
    
    if args.data:
        args.data = json.loads(args.data)
        
    if args.action == 'GET':
        result = api.get(args.entity, criteria=args.data)
    elif args.action =='POST':
        result =  api.post(args.entity, args.data)
        
    if args.pretty:
        print json.dumps(result, sort_keys=True, indent=4)
    else:
        print json.dumps(result)
        
if __name__ == '__main__':
    main()