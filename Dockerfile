FROM golang:alpine AS build

LABEL stage=builder

WORKDIR /build

ADD go.mod .

ADD go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/watcher ./cmd/watcher/watcher.go

FROM alpine

WORKDIR /app

COPY --from=build /app/watcher /app/watcher

CMD ["./watcher"]

