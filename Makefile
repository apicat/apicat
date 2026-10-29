.PHONY: build frontend test dev

frontend:
	cd web && npm ci && npm run build

build: frontend
	go build -o apicat-cli ./cmd/apicat-cli

test:
	go test ./...
	cd web && npm test

dev:
	cd web && npm run dev
