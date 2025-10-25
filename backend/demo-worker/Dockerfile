FROM golang:1.25 AS build
COPY . .
ENV GOPATH=""
ENV CGO_ENABLED=0
RUN go build -trimpath -a -ldflags="-w -s" ./cmd/demo-worker

FROM gcr.io/distroless/static
COPY --from=build /go/demo-worker /demo-worker
ENTRYPOINT ["/demo-worker"]
