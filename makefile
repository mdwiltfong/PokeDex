build:
	go build main.go
test:
	go test -v ./internal/tests/utils_test.go
run:
	./main