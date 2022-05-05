run:
	go run .\cmd\main.go

.PHONY: docker-clean
docker-clean:
	docker-compose down

.PHONY: docker-build
docker-build:
	docker-compose build server

.PHONY: docker-run
docker-run:	docker-clean docker-build 
	docker-compose up server

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
