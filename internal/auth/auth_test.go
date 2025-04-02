package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestPasswordCreate(t *testing.T) {
	password := "secret"
	hashedPass, err := HashPassword(password)
	if err != nil || hashedPass == "" {
		t.Errorf(fmt.Sprintf("pasword could not hash"))
	}
}

func TestPasswordValidate(t *testing.T) {
	password := "secret"
	hashedPass, err := HashPassword(password)
	if err != nil || hashedPass == "" {
		t.Errorf(fmt.Sprintf("pasword could not hash"))
		return
	}
	if err := CheckPasswordHash(hashedPass, password); err != nil {
		t.Errorf(fmt.Sprintf("could not unhash password %v\n", err))
		return
	}
}

func TestJWT(t *testing.T) {
	secret, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Errorf(fmt.Sprintf("Could not generate key secret %v", err))
		return
	}

	_, err = MakeJWT(uuid.New(), secret.X.String())
	if err != nil {
		t.Errorf(fmt.Sprintf("could not make JWT %v", err))
		return
	}
}

func TestJWTValidate(t *testing.T) {
	user := uuid.New()
	secret, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Errorf(fmt.Sprintf("Could not generate key secret %v", err))
		return
	}
	token, err := MakeJWT(user, secret.X.String())
	if err != nil {
		t.Errorf(fmt.Sprintf("could not make JWT %v", err))
		return
	}

	extractedUUID, err := ValidateJWT(token, secret.X.String())
	if err != nil {
		t.Errorf(fmt.Sprintf("could not validate JWT %v", err))
		return
	}
	if extractedUUID != user {
		t.Errorf("original uuid did not match jwt uuid")
		return
	}

}

func TestGetBearer(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer token")
	res, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("Token Extraction failed: %v", err)
		return
	}
	if res != "token" {
		t.Errorf("token was not properly processed")
		return
	}

}

func TestGetBearerBad(t *testing.T) {
	header := http.Header{}
	res, err := GetBearerToken(header)
	if err != nil {
		return
	}
	t.Errorf("nothing should be returned: %v", res)

}
