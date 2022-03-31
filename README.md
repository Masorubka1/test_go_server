1) make certificates for Os:

openssl req -new -newkey rsa:4096 -days 3650 -nodes -x509 -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=localhost" -addext "subjectAltName = DNS:localhost" -keyout server.key -out server.crt


2) upload the certificate to trust Arch : (sudo trust anchor --store path/to/crt/server.crt)
   debian based: https://unix.stackexchange.com/questions/90450/adding-a-self-signed-certificate-to-the-trusted-list
3) start the server locally with the command go run main.go
4) run tests to verify (if necessary)
5) use the server for its intended purpose
