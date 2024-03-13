package redis

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

var rdb *redis.Client

func InitRedisClient() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // Если требуется
		DB:       0,  // Номер БД в Redis, по умолчанию 0
	})
}

func GetRedisClient() *redis.Client {
	return rdb
}

type JWTHandlers struct {
	JwtPrivate *rsa.PrivateKey
	JwtPublic  *rsa.PublicKey
}

func NewAuthHandlers(jwtprivateFile string, jwtPublicFile string) *JWTHandlers {
	fmt.Fprintln(os.Stderr, "Creating JWT handlers")

	private, err := os.ReadFile(jwtprivateFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(private))
	public, err := os.ReadFile(jwtPublicFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	jwtPrivate, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		fmt.Println("Hello")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	jwtPublic, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return &JWTHandlers{
		JwtPrivate: jwtPrivate,
		JwtPublic:  jwtPublic,
	}
}

var h *JWTHandlers

func InitJWTHandlers() {
	privateFile := flag.String("private", "/auth_service/redis/signature/private_key.pem", "path to JWT private key `file`")
	publicFile := flag.String("public", "/auth_service/redis/signature/public_key.pem", "path to JWT public key `file`")
	flag.Parse()

	if privateFile == nil || *privateFile == "" {
		fmt.Fprintln(os.Stderr, "Please provide a path to JWT private key file")
		os.Exit(1)
	}

	if publicFile == nil || *publicFile == "" {
		fmt.Fprintln(os.Stderr, "Please provide a path to JWT public key file")
		os.Exit(1)
	}

	absoluteprivateFile, err := filepath.Abs(*privateFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	absolutePublicFile, err := filepath.Abs(*publicFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	h = NewAuthHandlers(absoluteprivateFile, absolutePublicFile)
}

func GetJWTHandlers() *JWTHandlers {
	return h
}
