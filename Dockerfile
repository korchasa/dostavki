FROM golang:1.13.7-alpine3.11 as build

WORKDIR /app

ENV \
    TERM=xterm-color \
    TIME_ZONE="Europe/Kiev" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOFLAGS="-mod=vendor" \
    GOLANGCI_VERSION="v1.23.3"

RUN \
    echo "## Prepare timezone" && \
    apk add --no-cache --update tzdata && \
    cp /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime && \
    echo "${TIME_ZONE}" > /etc/timezone && date && \
    echo "## Install golangci" && \
    wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_VERSION} && \
    golangci-lint --version
ADD . .
RUN \
    go env && \
    go version && \
    echo "  ## Test" && \
    go test ./... && \
    echo "  ## Lint" && \
    golangci-lint run ./... && \
    echo "  ## Build" && \
    go build -o app . && \
    echo "  ## Done"

FROM golang:1.13.7-alpine3.11
WORKDIR /app
COPY --from=build /app/app ./app
COPY --from=build /etc/localtime /etc/localtime
USER nobody:nobody
CMD ./app