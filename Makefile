all: tidy test testcov cuefmtcheck lintcue cuegen dirtycheck lintinsights

tidy:
	@echo "  >  Tidying go.mod ..."
	@go mod tidy
	@echo "  >  Tidying cue.mod ..."
	@cd schemas && cue mod tidy

test:
	@echo "  >  Running tests ..."
	@go vet ./...
	@go test ./...

testcov:
	@echo "Running tests and generating coverage output ..."
	@go test ./... -coverprofile coverage.out -covermode count
	@sleep 2 # Sleeping to allow for coverage.out file to get generated
	@echo "Current test coverage : $(shell go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+') %"

# Verify CUE formatting in ./schemas
cuefmtcheck:
	@echo "  >  Verifying CUE formatting in ./schemas ..."
	@cue fmt --check --files ./schemas

lint:
	@echo "  >  Linting Go files ..."
	@golangci-lint run

lintcue:
	@echo "  >  Linting CUE files (with module support) ..."
	@cd schemas && cue eval . --all-errors --verbose

cuegen:
	@echo "  >  Generating types from cue schema ..."
	@cd schemas && cue exp gengotypes .
	@mv schemas/cue_types_gen.go generated_types.go
	@go build -o cmd/types_tagger/types_tagger cmd/types_tagger/main.go
	@cmd/types_tagger/types_tagger generated_types.go
	@rm cmd/types_tagger/types_tagger

dirtycheck:
	@echo "  >  Checking for uncommitted changes ..."
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "  >  Uncommitted changes to generated files found!"; \
		echo "  >  Run make cuegen and commit the results."; \
		exit 1; \
	else \
		echo "  >  No uncommitted changes to generated files found."; \
	fi

oscalgenerate:
	@echo "  >  Generating OSCAL testdata from Gemara artifacts..."
	@mkdir -p artifacts
	@go run ./cmd/oscal_export catalog ./test-data/good-osps.yml --output ./artifacts/catalog.json
	@go run ./cmd/oscal_export guidance ./test-data/good-aigf.yaml --catalog-output ./artifacts/guidance.json --profile-output ./artifacts/profile.json

lintinsights:
	@echo "  >  Linting security-insights.yml ..."
	@curl -O --silent https://raw.githubusercontent.com/ossf/security-insights-spec/refs/tags/v2.1.0/schema.cue
	@cue vet -d '#SecurityInsights' security-insights.yml schema.cue
	@rm schema.cue
	@echo "  >  Linting security-insights.yml complete."

# Documentation site targets
CONTAINER_CMD := $(shell command -v podman 2> /dev/null || command -v docker 2> /dev/null)
VOLUME_FLAGS := $(shell [ "$$(uname -s)" = "Linux" ] && echo ":Z" || echo "")

check-container:
	@if [ -z "$(CONTAINER_CMD)" ]; then \
		echo "ERROR: Neither podman nor docker found."; \
		exit 1; \
	fi

serve: check-container
	@echo "  >  Starting Jekyll documentation site..."
	@echo "  >  Using container runtime: $(CONTAINER_CMD)"
	@$(CONTAINER_CMD) stop gemara-docs 2>/dev/null || true
	@$(CONTAINER_CMD) rm gemara-docs 2>/dev/null || true
	@echo "  >  Site will be available at: http://localhost:4000/gemara"
	@echo ""
	@$(CONTAINER_CMD) run --rm \
		--name gemara-docs \
		--volume="$$PWD/docs:/srv/jekyll$(VOLUME_FLAGS)" \
		--publish 4000:4000 \
		--publish 35729:35729 \
		docker.io/jekyll/jekyll:latest \
		jekyll serve --host 0.0.0.0 --livereload --force_polling

build: check-container
	@echo "  >  Building Jekyll documentation site..."
	@$(CONTAINER_CMD) run --rm \
		--volume="$$PWD/docs:/srv/jekyll$(VOLUME_FLAGS)" \
		docker.io/jekyll/jekyll:latest \
		jekyll build

clean: check-container
	@echo "  >  Cleaning generated files..."
	@rm -rf docs/_site docs/.jekyll-cache docs/.jekyll-metadata
	@echo "  >  Stopping and removing any running containers..."
	@$(CONTAINER_CMD) stop gemara-docs 2>/dev/null || true
	@$(CONTAINER_CMD) rm gemara-docs 2>/dev/null || true
	@echo "  >  Clean complete!"

stop: check-container
	@echo "  >  Stopping documentation server..."
	@$(CONTAINER_CMD) stop gemara-docs 2>/dev/null || true
	@$(CONTAINER_CMD) rm gemara-docs 2>/dev/null || true
	@echo "  >  Server stopped!"

restart: stop serve

.PHONY: tidy test testcov lintcue cuegen dirtycheck lintinsights serve build clean stop restart check-container
