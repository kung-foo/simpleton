VERSION := v0.0.1
BUILDSTRING := $(shell git log --pretty=format:'%h' -n 1)
VERSIONSTRING := simpleton version $(VERSION)+$(BUILDSTRING)

ifndef GOARCH
	GOARCH := $(shell go env GOARCH)
endif

ifndef GOOS
	GOOS := $(shell go env GOOS)
endif

OUTPUT := bin/simpleton-$(GOOS)-$(GOARCH)

ifeq ($(GOOS), windows)
	OUTPUT := $(OUTPUT).exe
endif

.PHONY: default gofmt all test clean goconvey docker-build

default: build

$(OUTPUT): main.go
	godep go build -v -o $(OUTPUT) -ldflags "-X main.VERSION \"$(VERSIONSTRING)\"" .
ifdef CALLING_UID
ifdef CALLING_GID
	@echo Reseting owner to $(CALLING_UID):$(CALLING_GID)
	chown $(CALLING_UID):$(CALLING_GID) $(OUTPUT)
endif
endif
	@echo
	@echo Built ./$(OUTPUT)

build: $(OUTPUT)

gofmt:
	gofmt -w .

update-godeps:
	rm -rf Godeps
	godep save

test:
	godep go test -cover -v ./...

clean:
	rm -f $(OUTPUT)
