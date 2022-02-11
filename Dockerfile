FROM golang:1.17-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

ENV GOARCH arm

RUN go mod download

RUN go build -o plex-init .

FROM alpine:3

COPY --from=builder /app/plex-init .

ENTRYPOINT [ "./plex-init" ]