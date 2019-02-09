install = dep ensure
lambda_env_vars = GOARCH=amd64 GOOS=linux

build:
	go build -i -ldflags="-s -w" -o bin/migrate-db ./migrate-db
	go build -i -ldflags="-s -w" -o bin/check ./check-for-updates
	go build -i -ldflags="-s -w" -o bin/shape-to-csv ./shape-to-csv
	go build -i -ldflags="-s -w" -o bin/import-csv-to-mysql ./import-csv-to-mysql

lambda:
	$(lambda_env_vars) go build -i -ldflags="-s -w" -o lambda-bin/shape-to-csv -tags lambda ./shape-to-csv
	$(lambda_env_vars) go build -i -ldflags="-s -w" -o lambda-bin/import-csv-to-mysql -tags lambda ./import-csv-to-mysql
	$(lambda_env_vars) go build -i -ldflags '-d -s -w' -o lambda-bin/migrate-db -tags lambda ./migrate-db

dry-deploy:
	serverless deploy --noDeploy

clean:
	rm -rf bin/*

test:
	go test ./import-csv-to-mysql
	go test ./migrate-db
