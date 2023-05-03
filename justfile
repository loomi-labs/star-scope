set dotenv-load

generate-models:
    cd api && go generate ./ent

generate-migrations:
    cd api && go run main.go database create-migrations

migrate:
    cd api && go run main.go database migrate

visualize-models:
    cd api && go run -mod=mod ariga.io/entviz ./ent/schema/