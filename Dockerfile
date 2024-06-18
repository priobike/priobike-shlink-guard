FROM golang
COPY main.go /go/src/main.go
RUN go build -o /go/bin/main /go/src/main.go
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD curl -f http://localhost:8000/health || exit 1
CMD ["/go/bin/main"]