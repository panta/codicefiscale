GO111MODULE=on
GOFLAGS= -mod=vendor
GO ?= go

GENERATED = comuni/comuni-generated-data.go

.PHONY: all
all: test $(GENERATED)

.PHONY: generate
$(GENERATED) generate: generate.go comuni/process-comuni.go
	$(GO) generate

.PHONY: test
test: $(GENERATED)
	$(GO) test

clean:
	-rm $(GENERATED)
