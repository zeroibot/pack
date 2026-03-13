go test -coverprofile=coverage.out ./dyn
go tool cover -html=coverage.out