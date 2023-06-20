# Ultimate Uptime Checker

### Running new migrations

```cmd
goose postgres "postgresql://postgres:password@localhost:5432/uptime_checker?sslmode=disable" create create_user_table sql
goose postgres "postgresql://postgres:password@localhost:5432/uptime_checker?sslmode=disable" up
```

### DB Codegen

```cmd
jet -source=postgresql -host=localhost -port=5432 -user=postgres -password=password -dbname=uptime_checker -schema=public -path=./schema -ignore-tables=goose_db_version,gue_jobs
```

### See Outdated

```cmd
go list -mod=mod -u -m -json all | go-mod-outdated -direct -update
```

### Update All Dependencies
```cmd
go get -u ./...
```
