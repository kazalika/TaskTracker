package jwt_handlers

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
)

const (
	privateKeyPath = "/auth_service/jwt/signature/private_key.pem"
	publicKeyPath  = "/auth_service/jwt/signature/public_key.pem"
)

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
	flag.Parse()

	absoluteprivateFile, err := filepath.Abs(privateKeyPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	absolutePublicFile, err := filepath.Abs(publicKeyPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	h = NewAuthHandlers(absoluteprivateFile, absolutePublicFile)
}

func GetJWTHandlers() *JWTHandlers {
	return h
}
