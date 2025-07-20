test:
	go test ./...

run:
	go run main.go https://monzo.com

build:
	go build -o monzo_web_crawler main.go