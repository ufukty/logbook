PSQL_USER = ufuktan

docker-build:
	docker build -t test .

migrate:
	cd database-migration && psql -U $(PSQL_USER) -d postgres -f create.sql