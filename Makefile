BINARY_NAME=fibonacci

try:
	go run ./pkg

build:
	go build -o ${BINARY_NAME} ./pkg/

run:
	./${BINARY_NAME}

register:
	docker run --network=host --rm ubercadence/cli:master -do test-domain domain register -rd 1

test:
	go test ./...

