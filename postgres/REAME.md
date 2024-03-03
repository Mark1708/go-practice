# Simple PostgreSQL projects

## Run apps
```bash
# Run project with github.com/lib/pq и database/sql
go run postgres/sql_pq.go postgres/todo_item.go

# Run project with github.com/jackc/pgx/v5 и github.com/jmoiron/sqlx
go run postgres/sqlx_pgx.go postgres/todo_item.go
```