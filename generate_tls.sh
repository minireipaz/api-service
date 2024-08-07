#
mkdir -p localcerts/ca localcerts/server localcerts/client

echo Generate the ca key
openssl genrsa -out localcerts/ca/key.pem 4096

echo Generate the ca certificate
openssl req -x509 -sha256 -new -nodes -key localcerts/ca/key.pem -days 3650 \
      -subj "/C=ES/ST=Baleares/L=Mallorca/O=Local Dev Org/OU=Org/CN=CA" \
      -extensions v3_ca \
      -out localcerts/ca/cert.pem


echo generating server certificate
openssl genrsa -out localcerts/server/key.pem 2048
openssl req -new \
      -subj "/C=ES/ST=Baleares/L=Mallorca/O=Local Dev Org/OU=Org/CN=localhost" \
      -key localcerts/server/key.pem \
      -out localcerts/server/signingReq.csr
openssl x509 -req -days 365 -in localcerts/server/signingReq.csr -CA localcerts/ca/cert.pem -CAkey localcerts/ca/key.pem -CAcreateserial -out localcerts/server/cert.pem
rm localcerts/server/signingReq.csr

openssl verify -CAfile localcerts/ca/cert.pem localcerts/server/cert.pem


echo generating client certificate
openssl genrsa -out localcerts/client/key.pem 2048
openssl req -new \
      -subj "/C=ES/ST=Baleares/L=Mallorca/O=Local Dev Org/OU=Org/CN=localhost" \
      -key localcerts/client/key.pem \
      -out localcerts/client/signingReq.csr
openssl x509 -req -days 365 -in localcerts/client/signingReq.csr -CA localcerts/ca/cert.pem -CAkey localcerts/ca/key.pem -CAcreateserial -out localcerts/client/cert.pem
rm localcerts/client/signingReq.csr

openssl verify -CAfile localcerts/ca/cert.pem localcerts/client/cert.pem

