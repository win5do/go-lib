.PHONY: test
test:
	go test -race ./...

bench:
	go test -run NONE  -bench . -benchtime 3s -cpu 2,4,8  .

lint:
	golangci-lint run -v
