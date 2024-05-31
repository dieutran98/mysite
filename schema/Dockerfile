FROM node:20-alpine

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
RUN echo "Installing system dependencies.." && \
      apk update && \
      apk add bash make openjdk17-jre && \
      yarn global add \
        @redocly/cli@1.8.1 \
        @openapitools/openapi-generator-cli && \
    echo "Cleanups.." && \
      rm -rf /var/cache/apk/*

ENV ROOT /schema

WORKDIR $ROOT
COPY Makefile package.json yarn.lock ./

RUN make install
