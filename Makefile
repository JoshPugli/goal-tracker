dev:
	cd backend && ~/go/bin/air

run:
	cd backend && go run main/server/main.go

build:
	cd backend && go build -o bin/server main/server/main.go