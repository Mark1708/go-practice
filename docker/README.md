# Build Docker Image

## Run Generator
```bash
# Installing generator
go install github.com/zeromicro/go-zero/tools/goctl@latest

# Generate
goctl docker -go gorilla_mux.go

# Build Image
docker build -t webscoket-app
```

## Resources
 - https://habr.com/ru/companies/otus/articles/660301/