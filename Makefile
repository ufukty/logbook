PSQL_USER = ufuktan

make:
	echo ""

build:
	docker build -t logbook-app-server .

build-with-logs:
	docker build --no-cache --progress plain -t logbook-app-server .

run: build
	docker run -it logbook-app-server

migrate:
	cd database-migration && psql -U $(PSQL_USER) -d postgres -f create.sql

psql-dev:
	psql -U ufuktan -d logbook_dev