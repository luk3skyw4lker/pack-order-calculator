export GOOSE_MIGRATION_DIR ?= ./src/database/migrations
export GOOSE_DRIVER ?= postgres
export GOOSE_DBSTRING ?= host=localhost port=5432 user=postgres password=postgres dbname=orders_db sslmode=disable

.SILENT:

run:
	if ! [ -f ./config/config.dev.yml ];    \
	then \
		go run ./src/main.go;    \
	else \
		go run ./stc/main.go --config-file=./config/config.dev.yml;    \
	fi

create_migration:
	goose create $(NAME) sql

migrate_up:
	goose up

migrate_up_by_one:
	goose up-by-one

migrate_up_to:
	goose up-to $(VERSION)

migrate_down:
	goose down

migrate_down_to:
	goose down-to $(VERSION)

db_status:
	goose status

db_reset:
	goose reset

watch:
	if ! [ -f ./config/config.dev.yml ];    \
	then \
		gow -c -v -r=false run ./src/main.go;    \
	else \
		gow -c -v -r=false run ./src/main.go --config-file=./config/config.dev.yml;    \
	fi

lint: fmt_docs
	golangci-lint run --allow-parallel-runners --fix

test:
	go test ./internal/... -race -coverprofile=fmt.coverage

generate_docs:
	cd src/ && swag init --parseDependency && cd ..

fmt_docs:
	swag fmt

cycle_containers: down_containers up_containers

down_containers:
	docker-compose down

up_containers:
	docker-compose up --build -d