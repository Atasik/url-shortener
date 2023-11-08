.PHONY:
.SILENT:

run:
    docker compose up link-shortener

run-db:
    docker compose up link-shortener-db
    
migrate:
    migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/urls?sslmode=disable' up        

test:
    go test -v ./...

swag:
    swag init -g internal/app/app.go