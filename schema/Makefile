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

bundle: validate
	@echo "Bundling YAML..."
	@npx swagger-cli bundle -o build/index.yml -t yaml ./src/index.yml

build: bundle
	@echo "Generating models and APIs..."
	@$(call generate, $(BUILD_DIRECTORY))
	@echo "Models and APIs generated successfully."
validate:
	@npx swagger-cli validate ./src/index.yml
