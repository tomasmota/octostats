FROM golang:1.17 as build-env

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /octostats

FROM gcr.io/distroless/base
COPY --from=build-env /octostats /
ENTRYPOINT ["/octostats"]