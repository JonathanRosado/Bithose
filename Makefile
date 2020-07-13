test_client:
	go build -o bithose cmd/bithose/main.go
	cd client && go test -v; \
	status=$$?; \
	rm ../bithose; \
	exit $$status

