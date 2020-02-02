FROM golang:1.13.4

RUN mkdir -p $GOPATH/src/github.com/asynccnu/table_service_v2
COPY . $GOPATH/src/github.com/asynccnu/table_service_v2
WORKDIR $GOPATH/src/github.com/asynccnu/table_service_v2

RUN go build -o main .
EXPOSE 8082
CMD ["/app/main"]