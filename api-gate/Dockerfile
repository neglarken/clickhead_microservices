FROM golang:1.20-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o main ./cmd/main.go

CMD [ "./main" ]