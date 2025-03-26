.PHONY: ssl
# Set the GOPROXY environment variable
export GOPROXY=https://goproxy.io,direct
export DEBUG=*

DOMAIN=localhost
IP=127.0.0.1
CN=Luoyy
BUILDTAGS=release

ssl:
	@echo 'authorityKeyIdentifier=keyid,issuer' > .v3.ext
	@echo 'basicConstraints=CA:FALSE' >> .v3.ext
	@echo 'keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment' >> .v3.ext
	@echo 'subjectAltName = @alt_names' >> .v3.ext
	@echo '[alt_names]' >> .v3.ext
	@echo 'IP.1 = ${IP}' >> .v3.ext
	@echo 'DNS.1 = ${DOMAIN}' >> .v3.ext

	@if [ ! -e root.key ]; then openssl genrsa -out root.key 2048; fi
	@if [ ! -e root.crt ]; then openssl req -x509 -new -nodes -key root.key -subj "/C=CN/CN=${CN}" -days 5000 -out root.crt; fi
	@if [ ! -e server.key ]; then openssl genrsa -out server.key 2048; fi
	@openssl req -new -key server.key -sha256 -subj "/CN=${DOMAIN}" -extensions v3_req -out server.csr
	@openssl x509 -req -in server.csr -CA root.crt -CAkey root.key -CAcreateserial -extfile .v3.ext -days 5000 -sha256 -out server.crt
	@rm -rf .v3.ext
