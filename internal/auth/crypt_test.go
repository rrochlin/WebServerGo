package auth

import (
	"fmt"
	"testing"
)

func TestPasswordCreate(t *testing.T) {
	password := "secret"
	hashedPass, err := HashPassword(password)
	if err != nil || hashedPass == "" {
		t.Errorf(fmt.Sprintf("pasword could not hash"))
		return
	}
}
