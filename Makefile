test:
	@go test -count=1 . -v

coverage:
	@go test -coverprofile=coverage.out . && go tool cover -html=coverage.out

format:
	@gofmt -s -w .