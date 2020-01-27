# ZROCK REST API 

### Social Platform for Musicians



## Migrations

https://github.com/golang-migrate/migrate


### Create 

    migrate create -ext  sql -dir migrations <migration_name>

### Up / Down

    migrate -path migrations -database postgres://postgres@localhost/database_name?sslmode=disable up/down
