build:
	dep ensure
	env GOOS=linux go build -i -ldflags="-s -w" -o bin/check check-for-updates/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
