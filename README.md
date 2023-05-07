# Star Scope

[Star Scope](https://star-scope.decrypto.online) brings clarity to the Cosmos ecosystem.

It is a webapp that gives Cosmonauts personalized notifications about their on-chain activity.

## How to run

To run the full app run the following command:
```bash
docker compose -f docker-full-app.yml up
```
Then go to http://localhost:8080.

How to clean up:
```bash
docker compose -f docker-full-app.yml down --volumes
```
