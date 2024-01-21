
build:
	env GOOS=linux go build -ldflags="-s -w" -o main cmd/app/lambda/main.go
	mkdir -p bin/
	zip bin/go-transaction-service.zip main
	rm main