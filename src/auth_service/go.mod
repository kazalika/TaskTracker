module auth_service

go 1.22.0

require (
	jwt_handlers v0.0.0
	kafka_handlers v0.0.0
	mongo_handlers v0.0.0
	task_service v0.0.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/gorilla/mux v1.8.1
	google.golang.org/grpc v1.63.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.47 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
)

require (
	github.com/gogo/protobuf v1.3.2
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace mongo_handlers => ./mongo_handlers

replace task_service => ../task_service

replace jwt_handlers => ./jwt_handlers

replace kafka_handlers => ./kafka_handlers
