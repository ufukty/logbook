PSQL_USER = ufuktan

migrate:
	cd database-migration && psql -U $(PSQL_USER) -d postgres -f create.sql