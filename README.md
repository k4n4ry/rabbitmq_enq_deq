# Abstract
enqueue process and dequeue process on rabbitMQ Server(at docker).  

# Usage

```
docker-compose up -d --build
go run cmd/deq/main.go
go run cmd/enq/main.go hello
```