BINARY  := ghisu
GOFLAGS := -trimpath

# Bump target: MAJOR, MINOR, or PATCH (case-insensitive)
BUMP ?= PATCH

.PHONY: all install build run test lint clean publish

all: build

install:
	go mod download

build:
	go build $(GOFLAGS) -o $(BINARY) .

run: build
	./$(BINARY)

test:
	go test ./...

lint:
	gofmt -w .
	go vet ./...

clean:
	rm -f $(BINARY)

# Compute the next semver tag from the latest git tag and BUMP.
# Usage:
#   make publish          # bumps patch (default)
#   make publish BUMP=MINOR
#   make publish BUMP=MAJOR
publish: lint test
	@LATEST=$$(git tag --list 'v*' --sort=-version:refname | head -1); \
	if [ -z "$$LATEST" ]; then LATEST="v0.0.0"; fi; \
	MAJOR=$$(echo $$LATEST | sed 's/^v//' | cut -d. -f1); \
	MINOR=$$(echo $$LATEST | sed 's/^v//' | cut -d. -f2); \
	PATCH=$$(echo $$LATEST | sed 's/^v//' | cut -d. -f3); \
	case "$$(echo $(BUMP) | tr '[:lower:]' '[:upper:]')" in \
	  MAJOR) MAJOR=$$((MAJOR+1)); MINOR=0; PATCH=0 ;; \
	  MINOR) MINOR=$$((MINOR+1)); PATCH=0 ;; \
	  PATCH) PATCH=$$((PATCH+1)) ;; \
	  *) echo "ERROR: BUMP must be MAJOR, MINOR, or PATCH (got '$(BUMP)')"; exit 1 ;; \
	esac; \
	NEXT="v$${MAJOR}.$${MINOR}.$${PATCH}"; \
	echo "Current: $${LATEST}  →  Next: $${NEXT}"; \
	git tag "$$NEXT" && \
	git push origin "$$NEXT" && \
	gh release create "$$NEXT" \
	  --title "$$NEXT" \
	  --generate-notes \
	  --verify-tag && \
	echo "Released $$NEXT"
