.DEFAULT_GOAL := install

BIN_FILE=4chan-thread-notifs

install:
	go build -o "${BIN_FILE}"

clean:
	go clean
	rm --force interview
	rm --force cp.out

test:
	go test

check:
	go test

cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out

