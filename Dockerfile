# build server
FROM golang:1.16 AS builder-server
COPY go.mod /src/
COPY go.sum /src/
RUN cd /src && go mod download
COPY . /src/
RUN cd /src && CGO_ENABLED=0 go build cmd/web/web.go

# runtime
FROM ineva/alpine:3.10.3
WORKDIR /app
COPY --from=builder-server /src/web /app
ENTRYPOINT /app/web