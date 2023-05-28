run:
	nodemon --exec go run main.go --signal SIGTERM
start-db:
	docker compose up -d
build:
	go build .