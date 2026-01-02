
.PHONY: help
help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo ""
	@echo "  Setup the app:"
	@echo "    setup-venv     - Setup venv (auto runs when needed)"
	@echo "    install-deps   - Install dependencies (auto runs when needed)"
	@echo "    build          - Build the application"
	@echo ""
	@echo "  Setup the environment:"
	@echo "    setup-dev      - Setup dev environment"
	@echo "    setup-prod     - Setup prod environment"
	@echo ""
	@echo "  Run the app:"
	@echo "    (Run an env setup command first, then run the app)"
	@echo "    init-new-story - Initialize a new story"
	@echo "    run-once       - Run the application once"
	@echo "    run-dry        - Run the application once, no external API calls"
	@echo "    run-daemon     - Run the application daily at configured time (see config/pipeline.yaml)"
	@echo ""
	@echo "  Test the app:"
	@echo "    test-all       - Run all tests"
	@echo "    test-unit      - Run all unit tests"
	@echo "    test-int       - Run all integration tests"
	@echo "    test-e2e       - Run all end-to-end tests"
	@echo ""
	@echo "  Misc:"
	@echo "    help           - Show this help message"
	@echo ""

# SETUP THE APP --------------------------------------------------------------

.PHONY: setup-venv
setup-venv:
	@echo "Venv setup not implemented yet"

.PHONY: install-deps
install-deps:
	go mod download
	@echo "Dependencies installed"

.PHONY: build
build:
	go build -o bin/storygen cmd/storygen/main.go
	@echo "Application built"
	
# SETUP THE ENVIRONMENT --------------------------------------------------------------

.PHONY: setup-dev
setup-dev:
	setup-venv
	source config/dev.env
	install-deps
	@echo "DEV environment setup complete"

.PHONY: setup-prod
setup-prod:
	setup-venv
	source config/prod.env
	install-deps
	@echo "PROD environment setup complete"

# RUN THE APP --------------------------------------------------------------

.PHONY: init-new-story
init-new-story:
	setup-venv
	./bin/storygen init
	@echo "New story initialized"

.PHONY: run-once
run-once:
	./bin/storygen run
	@echo "Application running once and will then terminate"

.PHONY: run-dry
run-dry:
	./bin/storygen run --dry-run
	@echo "Application running once and will then terminate, no external API calls"

.PHONY: run-daemon
run-daemon:
	./bin/storygen daemon
	@echo "Application running daily at configured time"

# TEST THE APP --------------------------------------------------------------

.PHONY: test-all
test-all:
	go test ./...
	@echo "All tests passed"

.PHONY: test-unit
test-unit:
	go test ./... -tags=unit
	@echo "Unit tests passed"

.PHONY: test-int
test-int:
	go test ./... -tags=integration
	@echo "Integration tests passed"

.PHONY: test-e2e
test-e2e:
	go test ./... -tags=e2e
	@echo "End-to-end tests passed"