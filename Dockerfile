FROM oven/bun:1 AS cssbase

WORKDIR /usr/src/app
COPY . .
RUN mkdir -p /css
RUN bunx tailwindcss -i ./css/index.css -o /css/output.css

FROM golang:1.22 AS build

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
# Build the Go application with static linking
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -ldflags '-w -s' -o /go/bin/app .

FROM alpine:latest AS app

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=cssbase /css/output.css /app/css/output.css
COPY --from=build /go/bin/app /app/server

EXPOSE 1323
USER 1000

CMD ["/app/server"]
