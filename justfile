set dotenv-load

generate-models:
    cd api && go generate ./ent

visualize-models:
    cd api && go run -mod=mod ariga.io/entviz ./ent/schema/

generate-migrations:
    cd api && go run main.go database create-migrations

migrate:
    cd api && go run main.go database migrate

generate-protobufs:
    protoc -I=proto/ --go_out=api/ --go_opt=module=github.com/shifty11/blocklog-backend \
            --go-grpc_out=api/ --go-grpc_opt=module=github.com/shifty11/blocklog-backend \
            proto/*.proto
