.PHONY: backend-run frontend-install frontend-run frontend-build test test-acceptance

backend-run:
	go run ./cmd/server

frontend-install:
	cd frontend && npm install

frontend-run:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

test:
	go test ./...

test-acceptance:
	go test ./backend/acceptance -v
