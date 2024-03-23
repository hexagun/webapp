FROM golang:1.22-alpine

RUN mkdir /app
WORKDIR /app

ADD server.go /app/
ADD go.mod /app/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /main
RUN ls -l
CMD ["/main"]