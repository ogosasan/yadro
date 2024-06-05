
build:
	go build -o web-server yadro/cmd/xkcd/

run:
	web-server

.PHONY: lint
lint:
	golangci-lint run

.PHONY: sec
sec:
	trivy fs .
	govulncheck

.PHONY: test
test:
	go test -coverprofile=coverage.out ./...

.PHONY: test1
test1:
	go tool cover -html=coverage.out -o coverage.html


.PHONY: e2e
e2e:
	./test.sh