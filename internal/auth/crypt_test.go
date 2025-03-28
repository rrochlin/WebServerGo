package auth

import (
	"testing"
)

func TestPasswordCreate(t *testing.T) {
	password := "secret"
	hashedPass, err := HashPassword(password)
	if err != nil || hashedPass == "" {
		t.Errorf("pasword could not hash")
		return
	}
}
