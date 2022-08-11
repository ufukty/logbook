PSQL_USER = ufuktan

make:
	echo ""

docker-build:
	docker build -t test .

migrate:
	cd database-migration && psql -U $(PSQL_USER) -d postgres -f create.sql