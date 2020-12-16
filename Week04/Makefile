TAG?=latest
SCHEMA?=User

dep:
	go get github.com/google/wire/cmd/wire
	go get github.com/go-sql-driver/mysql
	go get github.com/gorilla/mux
	go get github.com/swaggo/swag/cmd/swag
	go get github.com/swaggo/http-swagger

image:
	docker build -t week04-homework:$(TAG)

schema:
	ent init $(SCHEMA)

database:
	docker pull mysql:8.0.22
	docker run --name some-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:8.0.22

swagger:
	swag init -g cmd/server/main.go -o api