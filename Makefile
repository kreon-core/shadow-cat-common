.PHONY: dev-setup
dev-setup:
	go install github.com/go-delve/delve/cmd/dlv@latest							# for debugging
	go install golang.org/x/tools/cmd/goimports@latest							# for formatting imports
	go install github.com/daixiang0/gci@latest									# for organizing imports
	go install mvdan.cc/gofumpt@latest											# for formatting code
	go install github.com/segmentio/golines@latest								# for formatting long lines
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest	# for linting and formatting

.PHONY: pre-commit
pre-commit: install generate format lint clean

# ============================================================================================================

.PHONY: install
install:
	go mod download

.PHONY: generate
generate:

.PHONY: lint
lint:
	golangci-lint run || true

.PHONY: format
format:
	golangci-lint run --fix || true
	gofumpt -l -w -extra .
	goimports -w .
	gci write \
		--custom-order -s standard -s default -s "prefix(github.com/kreon-core/shadow-cat-common)" -s blank \
		--no-lex-order --skip-generated --skip-vendor .
	golines -w -m 120 .

.PHONY: build
build: generate

.PHONY: dev
dev: build
	go run .

.PHONY: clean
clean:
	go mod tidy
	go mod verify
	go mod edit -fmt
	go clean
	rm -rf build/
