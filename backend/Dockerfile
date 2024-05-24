FROM golang:1.22.3-alpine3.19 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

#COPY vendor ./vendor

#COPY api ./api
COPY cmd ./cmd
COPY store ./store
#COPY internal ./internal

COPY *.go ./

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/opengamifylms cmd/opengamifylms/main.go

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

EXPOSE 3000
EXPOSE 8080

COPY --from=build /app/bin/opengamifylms .

CMD ["/app/opengamifylms"]
