#!/bin/bash

openssl ecparam -name prime256v1 -genkey -noout -out ./certs/cakey.key
openssl req -x509 -new -nodes -key ./certs/cakey.key -subj "/CN=TestCA/C=MY" -days 3650 -out ./certs/cacert.pem
openssl ecparam -name prime256v1 -genkey -noout -out ./certs/server.key
openssl req -new -key ./certs/server.key -out ./certs/server.csr -config ./csr.conf
openssl x509 -req -in ./certs/server.csr -CA ./certs/cacert.pem -CAkey ./certs/cakey.key -CAcreateserial -out ./certs/server.pem -days 3650 -extfile ./csr.conf -extensions req_ext
openssl ecparam -name prime256v1 -genkey -noout -out ./certs/client.key
openssl req -new -key ./certs/client.key -out ./certs/client.csr -config ./csrclient.conf
openssl x509 -req -in ./certs/client.csr -CA ./certs/cacert.pem -CAkey ./certs/cakey.key -CAcreateserial -out ./certs/client.pem -days 3650 -extfile ./csrclient.conf -extensions req_ext
