# Simple Task Manager

## Description
Simple Task Manager is a RESTful API for managing tasks. It allows users to create, read, update, and delete tasks. The API is built using Go and uses MySQL as the database.

### How to run?
#### 1. Migrate tables
Make sure you have a tasks table
```
goose -dir ./db/migrations mysql "user:pass@tcp(localhost:3306)/tasks" up
```

#### 2. Configure environment variables
Use .env.example to create your own .env file.

#### 3. Generate Client Credentials
Save your generated client credential to request in the API
```
go run cmd/generator/main.go -name <client_name>
```

#### 4. Run server
```
go run cmd/server/main.go
```

#### 5. (Optional) Update swagger doc
For updates in routes or models, you need to regenerate the swagger documentation.
```
go install github.com/swaggo/swag/cmd/swag@latest
echo $GOPATH
export PATH=$PATH:$GOPATH/bin

swag init -g cmd/api/main.go
```

#### 6. Visit Swagger UI
Use your generated api key to use the API.
```
http://localhost/swagger/index.html
```
