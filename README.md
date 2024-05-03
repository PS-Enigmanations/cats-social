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
make migrateup #or

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

**Running k6**:

1. Ensure running `make dev` first
2. Run script:

```sh
make k6
```

**Docs**

1. Please install plugin `REST Client` at vscode
2. After create api, please create api documentation at folder `docs`. See example at `docs/auth.http`

**API**:

- http://localhost:8080/v1/user/register (`POST`)
- http://localhost:8080/v1/user/login (`POST`)
- http://localhost:8080/v1/cat (`GET`, `POST`)
- http://localhost:8080/v1/cat/match (`POST`)
- http://localhost:8080/v1/cat/match/approve
- http://localhost:8080/v1/cat/match/reject
