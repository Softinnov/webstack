webdev:
	npm install && node esbuild.js --dev

srvdev:
	docker-compose up -d db
	DIR=./ihm go run appTodo.go 

build-server:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/appTodo

build-client:
	cp ihm/index.html build/
	npm install && node esbuild.js

