# NOTE: No builder pattern apply as we need other tools for now
FROM golang:1.12.1-alpine3.9

# Set GOPATH environment variable and add gopath's bin to path
ENV GOPATH="/app/go"
ENV PATH="$PATH:$GOPATH/bin"
ENV APP_PATH="/src/genialo"

# Sets workdir and copy source code
WORKDIR $GOPATH/$APP_PATH
COPY . .

# Install dependencies & compile tool
RUN apk add --update build-base git \
    && make setup \
    && make build-container

# Sets workir and copy existing tests and gen command
COPY /app/go/src/gen/dist/genialo /usr/local/bin/genialo

# Runs gen command
CMD ["genialo"]
