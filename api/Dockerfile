FROM golang:1.22.0-alpine

# Install tools required for project
# Run `docker build --no-cache .` to update dependencies
RUN echo "Installing system dependencies.." && \
  apk add bash make && \
  echo "Cleanups.." && \
  rm -rf /var/cache/apk/*

ENV ROOT /api

WORKDIR $ROOT

COPY go.mod .
COPY go.sum .
COPY Makefile .
COPY .env .

RUN make install
