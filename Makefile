.PHONY: test build clean

name = fhir-server
build_dir = build
os = $(shell go env GOOS)
arch = $(shell go env GOARCH)
image_version = v1.0.0
image_tag = $(name)-image:$(image_version)
container = $(name)-container
air_hotreload = bin/air

export DOCKER_BUILDKIT=1

test:
	go test -v ./...

test-db:
	go test -timeout 30s -run ^TestNewDbHandler -v ./...

clean:
	rm -rf $(build_dir)/*
	go clean

run-server:
	$(build_dir)/$(name) --verbose

dev-server:
	if [ ! -f $(air_hotreload) ]; then \
		curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s; \
	fi

	$(air_hotreload) -c .air.toml

build-server:
	GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=0 go build -o $(build_dir)/$(name) .	
	
build-image:
	docker build -t $(image_tag) -f Dockerfile.multistage .

start-container:
	docker run -p $(port):$(port) --name $(container) --rm $(image_tag) 

stop-container:
	docker stop $(container)

fhir-examples:
	wget https://www.hl7.org/fhir/R4/examples-json.zip -O examples-json.zip
	unzip examples-json.zip -d examples-json
	rm -rf examples-json.zip