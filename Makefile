all:
	protoc --go_out=plugins:. api/services/uinames/uinames.proto
	protoc --go_out=plugins:. api/services/icndb/icndb.proto
	go build -o bin/uinames cmd/uinames/main.go
	go build -o bin/icndb cmd/icndb/main.go
	go build -o bin/comedian cmd/comedian/main.go
