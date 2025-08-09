.PHONY: dev run tidy build
#* API make commands
dev:
	cd backend && ~/go/bin/air

run:
	cd backend && go run /cmd/server/main.go

tidy:
	cd backend && go mod tidy

build:
	cd backend && go build -o bin/server cmd/server/main.go