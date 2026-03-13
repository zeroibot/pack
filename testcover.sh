go test -coverprofile=coverage.out ./list
go tool cover -html=coverage.out