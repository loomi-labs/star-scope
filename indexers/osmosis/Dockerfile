FROM golang:1.20-buster

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN rm Dockerfile

RUN go build -o /indexer

CMD [ "/indexer" ]
