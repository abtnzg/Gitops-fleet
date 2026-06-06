REGISTRY    ?= 851725313390.dkr.ecr.us-east-1.amazonaws.com
IMAGE       ?= apps/service-a
TAG         ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)

.PHONY: build push run test

build:
	docker build -t $(REGISTRY)/$(IMAGE):$(TAG) -t $(REGISTRY)/$(IMAGE):latest .

push: build
	aws ecr get-login-password --region us-east-1 | \
	  docker login --username AWS --password-stdin $(REGISTRY)
	docker push $(REGISTRY)/$(IMAGE):$(TAG)
	docker push $(REGISTRY)/$(IMAGE):latest

run:
	docker run --rm -p 8080:8080 $(REGISTRY)/$(IMAGE):latest

test:
	@echo "GET /"; curl -fsS http://localhost:8080/ | jq
	@echo "GET /healthz"; curl -fsS http://localhost:8080/healthz | jq
	@echo "GET /version"; curl -fsS http://localhost:8080/version | jq
	@echo "POST /echo"; curl -fsS -X POST -d 'hello platform' http://localhost:8080/echo | jq
