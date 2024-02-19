package rand

import (
	"crypto/rand"
	"encoding/base64"
)

// Nhan 1 integer n va tao slice do dai n, Read fills it with random byte
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

const RememberTokenBytes = 32

// RememberToken tao ra 1 remember token voi size 32
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}

func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}
