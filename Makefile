BINARY_NAME=fibonacci

try:
	go run .

build:
	go build -o ${BINARY_NAME} .

run:
	./${BINARY_NAME}

register:
	docker run --network=host --rm ubercadence/cli:master -do test-domain domain register -rd 1

test:
	go test ./...

lint:
	golangci-lint run
