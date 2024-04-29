# Cats Social

https://openidea-projectsprint.notion.site/Cats-Social-9e7639a6a68748c38c67f81d9ab3c769

**Database:**

```sh
DB_HOST=localhost
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=cats-social
DB_PORT=5432
DB_PARAMS="sslmode=disabled"
JWT_SECRET=a-very-secretive-secret-key
BCRYPT_SALT=8
```

**Run migration:**

1. Install [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)

2. Run scripts

```sh
make migrate #or

migrate -database "postgres://username:password@host:port/dbname?sslmode=disable" -path db/migrations up
```

**Setup:**

```sh
go mod download
```

**Running the server:**

```sh
make dev
```

**API**:

- http://localhost:8000/v1/cats (`GET`, `POST`)
