go run service.go --port=8081 --name=user_service
go run service.go --port=8082 --name=user_service
go run service.go --port=8083 --name=user_service


go run client.go --name=user_service