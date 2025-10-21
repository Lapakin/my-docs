DIR := $(CURDIR)
GOLINT_VERSION := v2.5.0

.PHONY:up
up:
	docker-compose up -d file-storage-minio
	$(MAKE) service-start SERVICE=user-management
	$(MAKE) service-start SERVICE=file-storage
	docker-compose up -d ui
	docker-compose up -d krakend

.PHONY:build
build:
	docker-compose build

.phony:stop
stop:
	docker-compose stop

.phony:down
down:
	docker-compose down

.PHONY:service-start
service-start:
	docker-compose up -d $(SERVICE)-postgres
	./hack/scripts/wait-for-postgres.sh ${SERVICE}-postgres
	$(MAKE) migrations-up SERVICE=${SERVICE}
	docker-compose up -d ${SERVICE}

.PHONY:service-build
service-build:
	docker-compose build ${SERVICE}

.PHONY:service-rebuild
service-rebuild:
	docker-compose stop ${SERVICE}
	docker-compose build ${SERVICE}
	$(MAKE) service-start SERVICE=${SERVICE}

.PHONY:service-stop
service-stop:
	docker-compose stop ${SERVICE}

.PHONY:migrations-up
migrations-up:
	@docker run --rm -v $(DIR)/pkg/${SERVICE}/repository/postgres/migrations:/migrations \
		--network file-storage_file-storage-network migrate/migrate -path=/migrations/ \
		-database "postgres://postgres:postgres@${SERVICE}-postgres:5432/$(subst -,_,${SERVICE})?sslmode=disable" up

.PHONY:migrations-down
migrations-down:
	@docker run --rm -v $(DIR)/pkg/${SERVICE}/repository/postgres/migrations:/migrations \
		--network file-storage_file-storage-network migrate/migrate -path=/migrations/ \
		-database "postgres://postgres:postgres@${SERVICE}-postgres:5432/$(subst -,_,${SERVICE})?sslmode=disable" down -all

.PHONY:generate-mocks
generate-mocks:
	./hack/scripts/generate-mocks.sh

.PHONY:go-unit-testing
go-unit-testing:
	./hack/scripts/run-go-unit-tests.sh

.PHONY:go-lint
go-lint:
	docker run -t --rm -v $(DIR):/app -v ~/.cache/golangci-lint/$(GOLINT_VERSION):/root/.cache -w /app golangci/golangci-lint:$(GOLINT_VERSION) golangci-lint run -v