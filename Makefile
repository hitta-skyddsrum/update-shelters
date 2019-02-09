install = dep ensure

build:
	$(install)
	go build -i -ldflags="-s -w" -o bin/check ./check-for-updates
	go build -i -ldflags="-s -w" -o bin/shape-to-csv ./shape-to-csv
	go build -i -ldflags="-s -w" -o bin/import-csv-to-mysql ./import-csv-to-mysql

lambda:
	$(install)
	GOARCH=amd64 GOOS=linux go build -i -ldflags="-s -w" -o lambda-bin/shape-to-csv -tags lambda ./shape-to-csv
	GOARCH=amd64 GOOS=linux go build -i -ldflags="-s -w" -o lambda-bin/import-csv-to-mysql -tags lambda ./import-csv-to-mysql

clean:
	rm -rf bin/*

test:
	go test ./import-csv-to-mysql
