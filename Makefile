all:
	protoc --go_out=plugins:. api/services/uinames/uinames.proto
	protoc --go_out=plugins:. api/services/icndb/icndb.proto
	CGO_ENABLED=0 GOOS=linux go build -a -o bin/uinames -tags netgo -ldflags '-w' cmd/uinames/main.go
	CGO_ENABLED=0 GOOS=linux go build -a -o bin/icndb -tags netgo -ldflags '-w' cmd/icndb/main.go
	CGO_ENBALED=0 GOOS=linux go build -a -o bin/comedian -tags netgo -ldflags '-w' cmd/comedian/main.go
