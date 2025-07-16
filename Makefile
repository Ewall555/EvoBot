NAME=evobot
BUILDDIR=build

BASEPATH := $(shell pwd)
BRANCH := $(shell git symbolic-ref HEAD 2>/dev/null | cut -d"/" -f 3)
BUILD := $(shell git rev-parse --short HEAD)

VERSION ?= $(BRANCH)-$(BUILD)
BuildTime:= $(shell date -u '+%Y-%m-%d %I:%M:%S%p')
VERSION ?= $(BRANCH)-$(BUILD)
TARGETARCH ?= amd64

UIDIR=frontend

.PHONY: docker
docker:
	@echo "build docker images"
	docker build --platform=linux/amd64 -f Dockerfile -t jumpserver-east/$(NAME):$(VERSION) .

.PHONY: clean
clean:
	rm -rf $(BUILDDIR)
	rm -rf $(UIDIR)/dist

update-version:
	@echo "▶ 版本号为: $(VERSION)"
	@cd installer && if [ -f .env ]; then \
		grep -v '^VERSION=' .env > .env.tmp; \
		echo "VERSION=$(VERSION)" >> .env.tmp; \
	else \
		echo "VERSION=$(VERSION)" > .env; \
	fi
	@echo "✓ env.txt 已更新"

.PHONY: package
package: update-version docker
	@echo "▶ 开始打包版本: $(VERSION)"
	@rm -rf $(BUILDDIR)
	@mkdir -p $(BUILDDIR)/$(NAME)/
	
	mv installer/.env.tmp $(BUILDDIR)/$(NAME)/.env
	cp installer/docker-compose.yml $(BUILDDIR)/$(NAME)/
	cp cmd/conf/config-example.yaml $(BUILDDIR)/$(NAME)/
	
	docker save -o $(BUILDDIR)/$(NAME)/$(NAME)-$(VERSION)-docker-image.tar jumpserver-east/$(NAME):$(VERSION)
	
	tar -czvf $(BUILDDIR)/$(NAME)-installer-$(VERSION)-package.tar.gz -C $(BUILDDIR)/$(NAME) .
	
	@echo "✓ 打包完成: $(BUILDDIR)/$(NAME)-installer-$(VERSION)-package.tar.gz"
