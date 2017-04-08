all:
	protoc --go_out=plugins:. api/services/uinames/uinames.proto
