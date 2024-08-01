.PHONY: install run api-doc modelgen codegen pgmodelgen

# Set default env as develop.
ifeq ($(ENV),)
ENV := develop
endif

-include .env.$(ENV).secret

WORK_DIR := .
API_DOC_DIR := $(WORK_DIR)/docs
API_DOC := $(API_DOC_DIR)/api.yaml

define _codegen
	oapi-codegen -generate chi-server -include-tags="$(1)" -package="$(1)" "$(API_DOC)"> "$(2)/$(1).gen.go"
endef

install:
	go mod download
	go install github.com/cosmtrek/air@v1.52.1
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0
	go install github.com/volatiletech/sqlboiler/v4@v4.16.0
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@v4.16.0

run:
	PORT=$(PORT) air -c .air.toml

# Copy api doc to be used with codegen.
api-doc:
	@echo "Copying API docs.."
	@mkdir -p $(API_DOC_DIR)
	@\cp -rf ./schema/build/index.yml $(API_DOC)
	@echo "Success!"

dtogen: api-doc
	@echo "Start dtogen..."
	$(eval ModelDir := $(WORK_DIR)/dtos)
	mkdir -p $(ModelDir)
	oapi-codegen -generate types -package="dtos" $(API_DOC) > "$(ModelDir)/dtos.gen.go"
	@echo "Success codegen! Check it out $(ModelDir)/dtos.gen.go"

codegen-help:
	@echo ""
	@echo "  Syntax:"
	@echo "  make codegen tag=<api tag>"
	@echo ""
	@echo "  Example:"
	@echo "  make codegen tag=health"
	@echo ""

codegen: api-doc
	@[ "${tag}" ] || ( make codegen-help ; exit 1 )
	@echo "Start codegen $(tag) ..."
	$(eval CODE_DIR := $(WORK_DIR)/features/$(tag))
	mkdir -p $(CODE_DIR)
	@$(call _codegen,$(tag),$(CODE_DIR))
	@echo "Successful codegen! Check it out at $(CODE_DIR)/$(tag).gen.go"

entitygen:
	@echo "Start entitygen..."
	PSQL_DBNAME=${POSTGRES_DB} PSQL_HOST=${POSTGRES_HOSTNAME} PSQL_USER=${POSTGRES_USER} PSQL_PASS=${POSTGRES_PASSWORD} PSQL_SSLMODE=${POSTGRES_SSL_MODE} sqlboiler psql --no-tests --wipe -o $(WORK_DIR)/entities -p entities
	@echo "Success dbmodelgen! Check it out $(WORK_DIR)/entities/*.go"

lint:
	golangci-lint run ./...
test:
	go test ./...
genmock:
	go generate ./...