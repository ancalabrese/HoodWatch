BINARY_NAME=HoodWatch
BINARY_PATH=bin/

build: 
	go build -o ${BINARY_PATH}/${BINARY_NAME} main.go

rund: build
	./${BINARY_PATH}/${BINARY_NAME} -d

run: build
	./${BINARY_PATH}/${BINARY_NAME}

clean:
	go clean
	rm -rf ${BINARY_PATH}