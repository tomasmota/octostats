run:
	go run main.go

build:
	go build -o bin/octostats

compile:
	@echo ">> Compiling for linux and windows"
	env GOOS=windows GOARCH=amd64 go build -o bin/octostats.exe
	env GOOS=linux GOARCH=amd64 go build -o bin/octostats

install:
	@echo ">> Installing octostats locally (~/go/bin/octostats)"
	go install .