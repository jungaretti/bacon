bin/bacon: */**/*.go
	go build -o bin/bacon

clean:
	rm -rf bin
