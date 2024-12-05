cvrg:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

dcup:
	docker-compose -f ./docker-compose.yml up

dcstop:
	docker-compose -f ./docker-compose.yml stop