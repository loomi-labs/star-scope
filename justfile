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
    cd client && trunk serve