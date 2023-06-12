.SILENT:
.DEFAULT_GOAL := run

run:
	docker-compose up --remove-orphans app

test:
	go test ./...