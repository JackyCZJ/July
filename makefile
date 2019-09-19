all:gotool
	go build .

build:
	go build .

clean:
	rm noghost
	find . -name "[._]*.s[a-w][a-z]" | xargs -i rm -f {}

gotool:
	gofmt -d .
	go vet . | grep -v vendor;true
ca:
	openssl req -new -nodes -x509 -out conf/server.crt -keyout conf/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=127.0.0.1/emailAddress=xxxxx@qq.com"

help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make test - test all file"
	@echo "make gotool - run go tool 'fmt' and 'vet'"
	@echo "make ca - generate ca files"
.PHONY: clean build gotool help ca