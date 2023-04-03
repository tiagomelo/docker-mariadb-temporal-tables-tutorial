# docker-mariadb-temporal-tables-tutorial

This is a tiny REST API to show how we can use temporal tables from [MariaDB](https://mariadb.org/).

Make sure you check `Postman` folder; I've exported a collection with invocation examples.

## running it

```
make run
```

## running unit tests

```
make test
```

## running integration tests

```
make int-test
```

## running linter

```
make lint
```

## generating api's documentation via Swagger

```
make swagger
```

## launching swagger ui

```
make swagger
```

Then head to `localhost`.

## references
- [go-swagger](https://github.com/go-swagger/go-swagger)
- [MariaDB](https://mariadb.org/)
- [System-versioned tables](https://mariadb.com/kb/en/system-versioned-tables/)

## related articles
- [running a dockerized linter](https://www.linkedin.com/pulse/golang-running-dockerized-linter-tiago-melo/)
- [declarative validations](https://www.linkedin.com/pulse/golang-declarative-validation-made-similar-ruby-rails-tiago-melo/)
- [database migrations](https://www.linkedin.com/pulse/go-database-migrations-made-easy-example-using-mysql-tiago-melo/)
