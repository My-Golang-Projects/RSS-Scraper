.PHONY: test clean all


run:
	go build && ./RSS-Scraper

# install that module
install_dot_env:
	go install github.com/joho/godotenv/cmd/godotenv@latest

db:
	docker run --name  rssagg -e POSTGRES_USER=pokemon -e POSTGRES_PASSWORD=secret123 -e POSTGRES_DB=rssagg -p 5432:5432 -d -d postgres:latest

up:
	cd sql/schema && \
	goose postgres postgres://pokemon:secret123@localhost:5432/rssagg up

down:
	cd sql/schema && \
	goose postgres postgres://pokemon:secret123@localhost:5432/rssagg down

sqlc_gen:
	sqlc generate

create_user:
	 curl -X POST -H 'Content-Type: application/json'  -d '{"name": "hari"}' http://localhost:8000/v1/users

# make API_KEY=<API_KEY>
create_feed:
	curl -X POST -H 'Content-Type: application/json' -H 'Authorization: ApiKey ${API_KEY}' -d '{"name": "Someone Blog", "url":"https://google.com"}'  http://localhost:8000/v1/feeds

get_feeds:
	curl http://localhost:8000/v1/feeds
