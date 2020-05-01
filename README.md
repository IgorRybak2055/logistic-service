# Ragger
The RAG system is a popular project management method for rating status reports

## Dependencies

For the work of Violets you will need:

* Postgresql version 12+
* Docker

For development, you will additionally need:

* Go 1.13+
* golang-migrate, `github.com / golang-migrate / migrate`

## Start

### With docker-compose
To run, you need database.env and ragger.env. The basis should be taken database_example.env 
and ragger_example.env rename them to database.env and ragger.env accordingly and 
make the necessary changes. Or simply:

```
make configs.env
```

After that you need to run docker-compose from the `deployments` directory with the command:

```
docker-compose  -f deployments/docker-compose.yml up
``` 

### With makefile 
To start using the makefile, check the rager's configs in the makefile and enter the command:
 
```
make run
```
  
Don't forget to start the database beforehand.


Check server availability:

```
curl localhost:8100/api/hello
``` 

## Migrations
A golang-migrate package, `github.com / golang-migrate / migrate`, is used to work with migration.
When the project starts, the database is automatically updated to the latest version.

For manual control using the command line:
```
migrate -source file://path/to/migrations -database postgres://localhost:5432/database (up or down) 2
```

## Swagger
A swaggo/swag package, `github.com/swaggo/swag`, is used to swagger documentation.

For generate swagger:
```
make swagger
```

## SQL
For SQL queries used formatter `http://www.dpriver.com/pp/sqlformat.htm`