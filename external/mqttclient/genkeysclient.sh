#!/bin/bash

openssl genrsa -out client.key 2048
openssl req -new -out client.csr -key client.key -extensions v3_ca -addext 'subjectAltName = IP:51.75.65.15' -subj '/C=US/ST=CA/L=SanFrancisco/O=RiotPot/OU=RND2/CN=vps-4d32dd92/'
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 360
chmod o=r client.key
cp client.key client.crt ca.crt /etc/mosquitto_client/certs/
