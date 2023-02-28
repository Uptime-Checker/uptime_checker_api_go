# Ultimate Uptime Checker

### Running new migrations

```cmd
goose postgres "postgresql://postgres:password@localhost:5432/uptime_checker_dev?sslmode=disable" create create_user_table sql
goose postgres "postgresql://postgres:password@localhost:5432/uptime_checker_dev?sslmode=disable" up
```

### DB Codegen

```cmd
jet -source=postgresql -host=localhost -port=5432 -user=postgres -password=password -dbname=uptime_checker -schema=public -path=./schema
```
