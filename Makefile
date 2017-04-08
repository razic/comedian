all:
	protoc --go_out=plugins:. api/services/uinames/uinames.proto
	go build -o bin/uinames cmd/uinames/main.go
