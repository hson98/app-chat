# Backend for AppChat

## About The Project
......
## Built With
* [Gin Web Framework](GoGin-url)
* [Generate SQL](Sqlc-url)
* [Driver Postgres](Pgx-url)
## Docs API
Sử dụng : [Docs swagger](https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format)</br>
Link docs: 
## Getting Started
### Prerequisites for develop
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

    ```bash
    brew install golang-migrate
    ```
- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [Gomock](https://github.com/uber-go/mock)

    ``` bash
    go install go.uber.org/mock/mockgen@latest
    ```
### How to run project ?
- Clone project
  ``` bash
  git clone https://github.com/hson98/app-chat.git
  ```
- Run Database, Redis
  ```base
   make dockerdb
  ```
- Create file app.env like example.env </br>
- Run db migration up all versions:
    ```bash
    make migrateup
    ```
- Run test:
    ```bash
    make test
    ```
- Run server:

    ```bash
    make server
    ```
## Run with docker
- Create file app.env like example.env </br>
 ```
  make docker
 ```
## Overview Architecture Project
![img_1.png](img_1.png)

## Overview Test Project
![img.png](img.png)
## Commands Note
- Gen ra docs api
  ```
  swag init -g cmd/api/main.go -o ./docs
  ```
- Create migration
  ```
  migrate create -ext sql -dir db/migration -seq init_schema
  ```
- Rebuild lại 1 service không bị cache
  ```
  docker-compose up -d --no-deps --build name_service
  ```
- Xóa cache build docker
  ```
  docker builder prune
  ```
- Xóa cache git
 ```
 git rm -r --cached .
 ```
