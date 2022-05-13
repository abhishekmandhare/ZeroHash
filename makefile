run:
	go run .\cmd\main.go --config .\config\local.yaml

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

.PHONY: lint
lint:
	 golangci-lint run

.PHONY: dep-install
dep-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.1