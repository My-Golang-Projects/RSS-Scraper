.PHONY: test clean all

run:
	go build && ./RSS-Scraper

# install that module
install_dot_env:
	go install github.com/joho/godotenv/cmd/godotenv@latest

db:
	docker run --name  rssagg -e POSTGRES_USER=pokemon -e POSTGRES_PASSWORD=secret123 -e POSTGRES_DB=rssagg -p 5432:5432 -d -d postgres:latest

migrate:
	cd sql/schema && \
	goose postgres postgres://pokemon:secret123@localhost:5432/rssagg up

sqlc_gen:
	sqlc generate
