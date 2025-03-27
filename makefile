run:
	go run ./cmd/app/main.go

compile:
	go build -o ./bin/app ./cmd/app/main.go

clean:
	rm -f ./bin/app

count:
	find . -type f -name "*.go" | xargs cat | wc -l
	find . -type f -name "*.go" -exec wc -l {} +