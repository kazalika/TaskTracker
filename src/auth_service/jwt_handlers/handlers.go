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
	privateKeyPath = "/auth_service/jwt_handlers/signature/private_key.pem"
	publicKeyPath  = "/auth_service/jwt_handlers/signature/public_key.pem"
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

func InitJWTHandlers() error {
	flag.Parse()

	absoluteprivateFile, err := filepath.Abs(privateKeyPath)
	if err != nil {
		return err
	}

	absolutePublicFile, err := filepath.Abs(publicKeyPath)
	if err != nil {
		return err
	}

	h = NewAuthHandlers(absoluteprivateFile, absolutePublicFile)
	return nil
}

func GetJWTHandlers() *JWTHandlers {
	return h
}
