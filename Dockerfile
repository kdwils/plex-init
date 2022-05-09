FROM golang as builder

ARG TARGETOS
ARG TARGETARCH
ENV GOARCH ${TARGETARCH:-amd64}
ENV GOOS ${TARGETOS:-linux}

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o plex-init .

FROM alpine
COPY --from=builder /app/plex-init .

ENTRYPOINT [ "./plex-init" ]