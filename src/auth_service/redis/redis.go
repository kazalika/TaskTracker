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
		Addr:     "localhost:6379",
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
	privateFile := flag.String("private", "/tmp/signature.pem", "path to JWT private key `file`")
	publicFile := flag.String("public", "/tmp/signature.pub", "path to JWT public key `file`")
	port := flag.Int("port", 8091, "http server port")

	flag.Parse()

	if port == nil {
		fmt.Fprintln(os.Stderr, "Port is required")
		os.Exit(1)
	}

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
