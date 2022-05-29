run:
	go build
	./dns-over-tor-resolver
docker:
	docker build . -t dnsserv
	docker rm -f dnsserv 2>/dev/null && docker run --name dnsserv -dit -p 127.0.0.1:53:5353/udp dnsserv
test:
	go clean -testcache ./...
	go test -v ./...