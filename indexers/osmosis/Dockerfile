FROM golang:1.20-buster

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN rm Dockerfile

COPY source_env.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/source_env.sh

RUN go build -o /indexer

CMD [ "/indexer" ]
