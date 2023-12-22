server := JishouSchedule
CONTAINER_NAME := dev-$(server)
IMAGE_NAME := $(shell echo $(server) | tr '[:upper:]' '[:lower:]')
EXPORT := 2000:2000
config := config.yaml

# Go编译器设置
GO := go
GOARCH_amd64 := amd64
GOARCH_arm64 := arm64
SYS_ARCH = $(shell uname -m)
ifeq ($(SYS_ARCH),arm64)
	SYS_ARCH := arm
endif

ifeq ($(SYS_ARCH),aarch64)
	SYS_ARCH := arm
endif

ifeq ($(SYS_ARCH),x86_64)
	SYS_ARCH := amd64
endif

# 输出目录
OUTPUT_DIR := bin
LINUX_OUTPUT := $(OUTPUT_DIR)

# 源文件
SOURCE := main.go

.PHONY: help clean run rm build docker

# all: linux

linux: help

build: linux_amd64 linux_arm64
	@echo "Build finished!"

linux_amd64:
	@echo "Compiling for Linux amd64..."
	@mkdir -p $(LINUX_OUTPUT)/linux_amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH_amd64) $(GO) build -o $(LINUX_OUTPUT)/linux_amd64/$(server) -ldflags "-w -s" $(SOURCE)
	@echo "Done"

linux_arm64:
	@echo "Compiling for Linux arm64..."
	@mkdir -p $(LINUX_OUTPUT)/linux_arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH_arm64) $(GO) build -o $(LINUX_OUTPUT)/linux_arm64/$(server) -ldflags "-w -s" $(SOURCE)
	@echo "Done"


clean:
	@echo "Cleaning..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Done"

docker: check-image
	@if docker ps -a --format "{{.Names}}" | grep $(CONTAINER_NAME); then \
		if docker ps -f "name=$(CONTAINER_NAME)" --format "{{.Status}}" | grep -q "Up"; then \
			echo "✅ Container $(CONTAINER_NAME) running"; \
		else \
			echo "🚀 Container $(CONTAINER_NAME)"; \
			docker start $(CONTAINER_NAME); \
		fi \
	else \
		echo "🚫 Container $(CONTAINER_NAME) not exist"; \
		echo "⏳ Container $(CONTAINER_NAME) starting"; \
		echo "Arch $(SYS_ARCH)"; \
		if [ "$(SYS_ARCH)" = "arm" ]; then \
			docker run -tid --restart=always -p $(EXPORT) \
			-v $(PWD)/$(LINUX_OUTPUT)/linux_arm64/$(server):/app/server \
			-v $(PWD)/$(config):/app/$(config) \
			-v $(PWD)/Deng.ttf:/app/Deng.ttf \
			--name $(CONTAINER_NAME) $(IMAGE_NAME); \
		else \
			docker run -ti --restart=always -p $(EXPORT) \
			-v $(PWD)/$(LINUX_OUTPUT)/linux_amd64/$(server):/app/server \
			-v $(PWD)/$(config):/app/$(config) \
			-v $(PWD)/Deng.ttf:/app/Deng.ttf \
			--name $(CONTAINER_NAME) $(IMAGE_NAME); \
		fi \
	fi

check-image:
	@if docker images --format "{{.Repository}}" | grep $(IMAGE_NAME); then \
		echo "🌟 Docker image $(IMAGE_NAME) exists"; \
	else \
		echo "📦 Docker image $(IMAGE_NAME) building"; \
		docker build -t $(IMAGE_NAME):latest .; \
	fi

rm:
	@echo "Will rm running docker container."
	docker rm -f $(CONTAINER_NAME)

run:
	@echo "🚀 Start Application $(server)"
	chmod +x $(PWD)/$(LINUX_OUTPUT)/linux_amd64/$(server)
	$(PWD)/$(LINUX_OUTPUT)/linux_amd64/$(server)

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  Build         - Build for all platforms"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker        - Build Docker image and run container"
	@echo "  rm            - Remove Docker running container"
	@echo "  run           - Start app server"
