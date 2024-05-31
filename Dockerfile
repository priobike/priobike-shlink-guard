FROM golang
COPY main.go /go/src/main.go
RUN go build -o /go/bin/main /go/src/main.go
CMD ["/go/bin/main"]