.PHONY: install lint bundle build

# Set default env as develop.
ifeq ($(ENV),)
ENV := develop
endif

# Build directory
BUILD_DIRECTORY := ./build

define generate
	npx openapi-generator-cli generate -i $(1)/index.yaml -g typescript-axios -o $(1)/typescript --strict-spec=true --reserved-words-mappings "class=class" --additional-properties=modelPackage=models,apiPackage=apis,withSeparateModelsAndApi=true
endef

# Install related tool
install:
	@echo "Installing dependencies..."
	@yarn install --frozen-lockfile
	@echo "Dependencies installed successfully."

lint:
	@echo "Linting YAML..."
	@npx @redocly/cli lint src/index.yml --skip-rule=info-license

bundle:
	@echo "Bundling YAML..."
	@npx @redocly/cli bundle src/index.yml --output $(BUILD_DIRECTORY)/index.yaml --ext yaml

build: bundle
	@echo "Generating models and APIs..."
	@$(call generate, $(BUILD_DIRECTORY))
	@echo "Models and APIs generated successfully."
gen:
	@npx openapi-generator-cli generate -i ./src/index.yml -g ruby -o /tmp/test/
validate:
	@npx swagger-cli validate ./src/index.yml
swagger-bundle:validate
	@npx swagger-cli bundle -o build/index.yml -t yaml ./src/index.yml