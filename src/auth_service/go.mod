module auth_service

require (
	api_handlers v0.0.0-00010101000000-000000000000
	jwt_handlers v0.0.0
	redis v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
)

replace api_handlers => ./go

replace redis => ./redis

replace jwt_handlers => ./jwt

go 1.18
