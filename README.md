# Star Scope
## Description


## How to run

To run the full app run the following command:
```bash
docker compose -f docker-full-app.yml up
```

How to clean up:
```bash
docker compose -f docker-full-app.yml down --volumes
```


## Architecture

### Frontend

The frontend is a rust application that uses the [Sycamore](https://sycamore-rs.netlify.app/) framework. 
It is compiled to wasm and served by the caddy server. I use [tailwindcss](https://tailwindcss.com/) for styling.

### Server

#### gRPC server
The gRPC server is a golang application with a postgres database. It is responsible for storing the data and serving it to the frontend.

#### event consumer
The event consumer is a golang application that listens to a kafka topic for new blockchain transactions and stores them in the database.

### Indexers

Right now there is only one indexer. I plan to add more and cover the whole Cosmos ecosystem.

#### Osmosis indexer
The Osmosis indexer is a golang application that listens to new blocks on the osmosis blockchain and publishes them to a kafka topic.
