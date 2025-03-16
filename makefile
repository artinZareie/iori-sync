run:
	go run ./cmd/app/main.go

compile:
	go build -o ./bin/app ./cmd/app/main.go

clean:
	rm -f ./bin/app
