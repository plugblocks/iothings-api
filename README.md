# things-API

Things API

## How to install

Run openssl genrsa -out base.rsa 1024 to generate the private key. Store this key at the root level.

Run openssl rsa -in base.rsa -pubout > base.rsa.pub to generate the public key. Store this key at the root level.
