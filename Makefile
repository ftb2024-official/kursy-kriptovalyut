cvrg:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

dcup:
	docker-compose -f ./docker-compose.yaml up

dcstop:
	docker-compose -f ./docker-compose.yaml stop