# go-currency-history

This repository is for my article on Medium about Web Scrapping and Concurrency in Go.

# How to run the project

Run the server
```shell
go run ./cmd/main.go
```

Run with docker-compose
```shell
docker-compose up
```

To test the APIs open:

[http://localhost:8080/v1/currency](http://localhost:8080/v1/currency)

and

[http://localhost:8080/v1/currency/history?from=IDR&to=JPY&start_date=2022-03-01&end_date=2022-03-19](http://localhost:8080/v1/currency/history?from=IDR&to=JPY&start_date=2022-03-01&end_date=2022-03-19)
