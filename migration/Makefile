.PHONY: install migrate down drop


MIGRATE := migrate -source file://ddl -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOSTNAME):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL_MODE)'

install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

migrate:
	@echo "Migrating..."
	@if ! $(MIGRATE) up; then \
		echo "Error: Migration failed"; \
		exit 1; \
	fi
	@echo "Done!"

down:
	@echo "Migrating(down) DB.."
	@echo y | $(MIGRATE) down 1
	@echo "Done!"

drop:
	@echo "Dropping..."
	@if ! $(MIGRATE) down; then \
		echo "Error: Dropping migration failed"; \
		exit 1; \
	fi
	@echo "Done!"

drop-to-version:
	@[ "${version}" ] || ( echo "please input version" ; exit 1 )
	@echo "Dropping to $(version)..."
	@if ! $(MIGRATE) force $(version); then \
		echo "Error: Dropping migration failed"; \
		exit 1; \
	fi
	@echo "Done!"
