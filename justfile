set dotenv-load

generate-models:
    cd server && go generate ./ent

generate-migrations:
    cd server && go run main.go database create-migrations

migrate:
    cd server && go run main.go database migrate

visualize-models:
    cd server && go run -mod=mod ariga.io/entviz ./ent/schema/

client:
    cd client && export GRPC_WEB_ENDPOINT_URL=http://127.0.0.1:8090 && export DISCORD_CLIENT_ID=953923165808107540 && export WEB_APP_URL=http://test.mydomain.com:8080 && trunk serve