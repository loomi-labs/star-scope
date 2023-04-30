# Blocklog-backend
## Description


## How to run

Setup the environment variables:
```bash
cp api/.env.template api/.env
```

To run the Postrgres database and the gRPC server, run the following command:
```bash
docker compose -f docker-full-app.yml up
```