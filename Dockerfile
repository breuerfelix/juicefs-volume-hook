FROM golang:1.24 as builder

WORKDIR /workspace
COPY go.mod go.mod
RUN go mod download

COPY *.go .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/app .
USER 65532:65532

ENTRYPOINT ["/manager"]
