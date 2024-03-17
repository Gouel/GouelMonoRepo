openssl genpkey -algorithm RSA -out priv.pem -pkeyopt rsa_keygen_bits:2048
openssl req -new -key priv.pem -out csr.pem -config openssl.cnf
openssl x509 -req -in csr.pem -signkey priv.pem -out cert.pem -days 365 -extfile openssl.cnf -extensions req_ext
