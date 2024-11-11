build:
	go build -o bin/postgres-benchmark-tool ./cmd/app

docker-build:
	docker-compose build app

docker-run:
	docker-compose up

run-app:
	docker-compose up app

clean:
	rm -rf bin/ postgres-benchmark-tool
	docker-compose down

deps:
	go mod tidy

reset:
	docker-compose down -v
	docker-compose up
