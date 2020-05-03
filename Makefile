export HTTP_HOST := localhost
export HTTP_ADDR := :8888
export LOG_LEVEL := debug

export LOCAL_DEBUG := false

export DATABASE_HOST := localhost
export DATABASE_NAME := logistic
export DATABASE_USER := logistic
export DATABASE_PASSWORD := logistic
export DATABASE_SSLMODE := disable

export MIGRATIONS_FILES := /migrations

export TOKEN_PASSWORD := password_for_signing_token

export SMTP_URL := smtp.gmail.com
export SMTP_PORT := 587
export SMTP_SENDER := email_address
export SMTP_SENDER_PASS := email_password
export SMTP_TEMPLATE := ./pkg/email

export GOPROXY=direct

.PHONY: configs.env
configs.env:
	cp config/logistic_example.env config/logistic.env
	cp config/database_example.env config/database.env

.PHONY: test
test:
	go test -v -race ./...

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: swagger
swagger:
	$(HOME)/go/bin/swag init -g ./cmd/main.go

export POSTGRESQL_URL='postgres://logistic:logistic@localhost:5432/logistic?sslmode=disable'

.PHONY: migrate
migrate:
	migrate -database ${POSTGRESQL_URL} -path migrations up
