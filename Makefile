build:
	go build -o bin/bacon .

test-unit:
	go test ./...

test-system: build
	git submodule update --init
	./test/bats/bin/bats test/test.bats

clean:
	rm -rf bin
