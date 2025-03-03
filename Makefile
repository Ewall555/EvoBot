BUILDDIR=build

BASEPATH := $(shell pwd)
BRANCH := $(shell git symbolic-ref HEAD 2>/dev/null | cut -d"/" -f 3)
BUILD := $(shell git rev-parse --short HEAD)
KOKOSRCFILE := $(BASEPATH)/cmd/koko/
KUBECTLFILE := $(BASEPATH)/cmd/kubectl/
HELMFILE := $(BASEPATH)/cmd/helm/

VERSION ?= $(BRANCH)-$(BUILD)
TARGETARCH ?= amd64

UIDIR=frontend
NPMINSTALL=yarn install
NPMBUILD=yarn build


evo-ui:
	@echo "build ui"
	@cd $(UIDIR) && $(NPMINSTALL) && $(NPMBUILD)


.PHONY: docker
docker:
	@echo "build docker images"
	docker buildx build --build-arg VERSION=$(VERSION) -t jumpserver-east/evo:$(VERSION)  .

.PHONY: docker
package:
	mkdir -p build
	cp installer/* build/
	cp cmd/server/conf/config-example.yaml build/
	docker save -o  build/evo-$(VERSION)-docker-image.tar jumpserver-east/evo:$(VERSION)
	tar -czvf build/evo-$(VERSION)-package.tar.gz build/*
