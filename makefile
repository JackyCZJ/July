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

help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make test - test all file"
	@echo "make gotool - run go tool 'fmt' and 'vet'"
.PHONY: clean build gotool help