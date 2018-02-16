build:
	dep ensure
	env GOOS=linux go build -i -ldflags="-s -w" -o bin/check check-for-updates/*.go
	env GOOS=linux go build -i -ldflags="-s -w" -o bin/shape-to-geojson shape-to-geojson/*.go
