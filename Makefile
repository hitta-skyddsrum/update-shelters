install = dep ensure
folders = check-for-updates import-csv-to-mysql migrate-db rest-api shape-to-csv
lambda_env_vars = GOARCH=amd64 GOOS=linux

build:
	$(foreach folder,$(folders),go build -i -ldflags="-s -w" -o bin/$(folder) ./$(folder);)

lambda:
	$(lambda_env_vars) go build -i -ldflags="-s -w" -o lambda-bin/shape-to-csv -tags lambda ./shape-to-csv
	$(lambda_env_vars) go build -i -ldflags="-s -w" -o lambda-bin/import-csv-to-mysql -tags lambda ./import-csv-to-mysql
	$(lambda_env_vars) go build -i -ldflags="-s -w" -o lambda-bin/migrate-db -tags lambda ./migrate-db

dry-deploy:
	serverless deploy --noDeploy

clean:
	rm -rf bin/*

test:
	go test ./import-csv-to-mysql
	go test ./migrate-db
	go test ./rest-api

coverage:
	./tools/coverage.sh
