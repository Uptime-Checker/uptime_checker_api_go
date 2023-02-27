# Ultimate Uptime Checker

### Running new migrations

```cmd
goose postgres "postgresql://postgres:password@localhost:5432/uptime_checker_dev?sslmode=disable" create create_user_table sql
goose postgres "postgresql://postgres:password@localhost:5432/uptime_checker_dev?sslmode=disable" up
```
