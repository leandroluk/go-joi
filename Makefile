GO          ?= go
PKG         := ./...
COVERPKG    := ./joi
TMPDIR      := ./.tmp
COVERFILE   := $(TMPDIR)/coverage.out
COVERHTML   := $(TMPDIR)/coverage.html

# Cross-platform commands for creating/removing files
ifeq ($(OS),Windows_NT)
    MKDIR = if not exist $(TMPDIR) mkdir $(TMPDIR)
    RMRF  = del /Q
else
    MKDIR = mkdir -p $(TMPDIR)
    RMRF  = rm -f
endif

.PHONY: test test.ci test.func test.html coverage.save fmt vet lint clean help

## Run all tests (no coverage)
test:
	$(GO) test $(PKG)

## Run tests with coverage (generate $(COVERFILE))
test.ci:
	$(MKDIR)
	$(GO) test $(PKG) -coverpkg=$(COVERPKG) -coverprofile=$(COVERFILE) -covermode=atomic

## Open the coverage report in the browser (CROSS-PLATFORM)
#  Note: does NOT generate a file, just opens directly.
test.html: test.ci
	$(GO) tool cover -html=$(COVERFILE)

## Show coverage by function in the terminal (aggregate view)
test.func: test.ci
	$(GO) tool cover -func=$(COVERFILE)

## Save the HTML coverage report to a file (does not open automatically)
coverage.save: test.ci
	$(GO) tool cover -html=$(COVERFILE) -o $(COVERHTML)
	@echo "Report saved at: $(COVERHTML)"

## Formatters/Linters
fmt:
	$(GO) fmt $(PKG)

vet:
	$(GO) vet $(PKG)

lint: fmt vet

## Clean up
clean:
	$(RMRF) $(COVERFILE) $(COVERHTML)

## Help (simple and portable)
help:
	@echo "Targets:"
	@echo "  make test           - Run tests"
	@echo "  make test.ci        - Run tests with coverage (profile)"
	@echo "  make test.html      - Open coverage report in browser (cross-platform)"
	@echo "  make test.func      - Show coverage by function in terminal"
	@echo "  make coverage.save  - Save HTML report to $(COVERHTML)"
	@echo "  make fmt vet lint   - Formatters/Linters"
	@echo "  make clean          - Remove coverage files"
