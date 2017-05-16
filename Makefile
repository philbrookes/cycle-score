SHELL = /bin/bash
VERSION = 1.0.0
NAME = cycle-score.com

.PHONY: build_and_serve
build_and_serve: build serve

.PHONY: build
build: build_server build_app

.PHONY: build_server
build_server:
	cd cmd/rest-api && go build .

.PHONY: build_app
build_app:
	cd web && ng build --output-path public

.PHONY: serve
serve:
	./cmd/rest-api/rest-api