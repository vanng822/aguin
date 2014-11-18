
Encrypt:

1. AES encrypt data with 256-bit key, CFB mode and autopadding
2. Signed it with hmac/sha256
3. Urlencoded base64 on both data/message and expected signature
4. Concat the message in format: base64_expected_signature.base64_encrypted_message

Decrypt:

1. Reverse of Encrypt


Base64 urlencode:

1. + => -
2. / => _

Base64 urldecode:

1. - => +
2. _ => /

API:

1. GET: Getting records for an certain Entity
2. POST: creating new Entity entry. Entity should be in json format and sent as parameter message.

All API calls must make with a header X-AGUIN-API-KEY and an optional header X-AGUIN-API-SECRET. Secret header should only send if you have https.


Entity:

1. entity: string
2. data: json data with primitive type. Integer, float, boolean and array of integer or float.


Test data:

1. input: WmtlMxl4d_VTfUYnl-A0Uycpr2e3VswKDwoPd03XtoY=.JZMk6FZloYh5BL0K7dHGSyTqB4lTgd9annrFEgLTELnxR3bHweL2

2. decrypted: {"entity":"something2"}