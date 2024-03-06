FROM oven/bun:1 as base
WORKDIR /usr/src/app

COPY . .

RUN bunx tailwindcss -i ./css/index.css -o ./css/output.css

FROM golang:1.22 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate

COPY --from=base /usr/src/app/css/output.css /usr/src/app/css/output.css

# Build the Go application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/app /app

CMD ["/app"]
