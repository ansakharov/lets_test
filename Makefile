PORT?=9999
APP_NAME?=test-app

clean:
	rm -f ${APP_NAME}

build: clean
	go build -o ${APP_NAME}

run: build
	PORT=${PORT} ./${APP_NAME}

test:
	go test -v -count=1 ./...

test100:
	go test -v -count=100 ./...

race:
	go test -v -race -count=1 ./...

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: gen
gen:
	mockgen -source=internal/pkg/repository/order/repository.go \
	-destination=internal/pkg/repository/order/mocks/mock_repository.go