# RSS-Scraper

## Database

- [sqlc](https://github.com/sqlc-dev/sqlc) Generate type-safe code from SQL
- `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

```sql
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4)
RETURNING *;
```

1. I want a new function called `CreateUser` that return one user
2. `$n` create the function with n parameters
3. `RETURNING *` return that user we just inserted

### Migration Tool

- [Goose](https://github.com/pressly/goose)
- `go install github.com/pressly/goose/v3/cmd/goose@latest`

- works with the sql comment

```bash
cd sql/schema
goose postgres <conn_string> up
```

- it will create a `goose_db_version` which is managed by goose
