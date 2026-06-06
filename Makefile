.PHONY: all build run dev clean frontend

# Default: build everything and run
all: build run

# Build the single binary (with embedded frontend)
build: frontend
	go build -o work-manager .

# Run the server
run:
	./work-manager

# Build and run in one command
start: build run

# Build frontend only
frontend:
	cd frontend && npm install --silent && npm run build

# Dev mode: run Go backend + Vite dev server separately
dev:
	@echo "Start Go backend in one terminal: go run ."
	@echo "Start Vite dev in another: cd frontend && npm run dev"
	@echo "Then open: http://localhost:5173 (Vite proxies /api to :8080)"

# Clean build artifacts
clean:
	rm -f work-manager work_manager.db
	rm -rf frontend/dist
