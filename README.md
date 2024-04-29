# Cats Social

https://openidea-projectsprint.notion.site/Cats-Social-9e7639a6a68748c38c67f81d9ab3c769

**Database:**

```sh
SECRET_KEY=a-very-secretive-secret-key
DATABASE_HOST=localhost
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=cats-social
DATABASE_PORT=5432
```

**Run migration:**

1. Install [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)

2. Run scripts

```sh
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
