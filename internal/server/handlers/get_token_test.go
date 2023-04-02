package handlers

import (
	"testing"
	"time"
)

func Test_generateTokenFromHeader(t *testing.T) {
	header, err := createHeader("127.0.0.1", 5, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	token := generateTokenFromHeader(header)
	if len(token) == 0 {
		t.Fatal("token is empty")
	}
	t.Logf("token: %s", token)
}
