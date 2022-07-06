.PHONY: vendor
vendor:
	go mod vendor

build:
	go build

lambda-build:
	GOOS=linux go build -mod=vendor -o /asset-output/main