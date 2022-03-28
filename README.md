1) сделать сертификаты для Os:
openssl req \
         -new \
         -newkey rsa:4096 \
         -days 3650 \
         -nodes \
         -x509 \
         -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=localhost" \
         -keyout server.key \
         -out server.crt
2) закинуть сертификат в trust множество Arch : (sudo trust anchor --store path/to/crt/server.crt)
   debian based: https://unix.stackexchange.com/questions/90450/adding-a-self-signed-certificate-to-the-trusted-list
3) запустить локально сервер командой go run main.go
*4) запустить тесты для проверки
5) воспользоваться сервером по назначению