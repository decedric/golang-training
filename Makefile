BINARY_NAME=fibonacci

try:
	go run ./pkg

build:
	go build -o ${BINARY_NAME} ./pkg/

run:
	./${BINARY_NAME}


test:
	go test ./...

