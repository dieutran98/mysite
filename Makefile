.PHONY: up build build-one down

# Set default env as develop if ENV is not specified.
ENV ?= develop

# Alias command for docker-compose.
COMPOSE := ENV=$(ENV) COMPOSE_HTTP_TIMEOUT=300 docker-compose

# Boot all containers
up:
	$(COMPOSE) up -d

# Build all images
build:
	$(COMPOSE) build

# Build one image with a specific tag
build-one:
	$(COMPOSE) build $(tag)

# Kill all containers
down:
	$(COMPOSE) kill
