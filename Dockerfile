FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go mod tidy
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/scores-app ./main.go

FROM alpine:3.13
RUN apk --no-cache add bash curl ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
ENTRYPOINT /go/bin/scores-app