FROM golang as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

ENV GOARCH ${TARGETARCH:-amd64}
ENV GOOS ${TARGETOS:-linux}

COPY go.* ./
RUN go mod download

COPY . ./

RUN go mod download

RUN go build -o plex-init .

FROM alpine

COPY --from=builder /app/plex-init .

ENTRYPOINT [ "./plex-init" ]