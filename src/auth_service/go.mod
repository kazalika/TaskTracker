module auth

require (
	auth/redis v0.0.0-00010101000000-000000000000
	sw v1.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
)

replace sw v1.0.0 => ./go

replace auth/redis => ./redis

go 1.18
