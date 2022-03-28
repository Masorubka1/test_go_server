1) сделать сертификаты для Os:
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
2) закинуть сертификат в trust множество Arch : (sudo trust anchor --store path/to/crt/server.crt)
   debian based: https://unix.stackexchange.com/questions/90450/adding-a-self-signed-certificate-to-the-trusted-list
3) запустить локально сервер командой go run main.go
4*) запустить тесты для проверки
5) воспользоваться сервером по назначению