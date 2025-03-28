package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

func Hash(key string) string {
	hasher := sha256.New()

	hasher.Write([]byte(key))

	hashSum := hasher.Sum(nil)

	return hex.EncodeToString(hashSum)
}

func verifyPassword(storedHash, providedPassword string) bool {
	return storedHash == providedPassword
}

func RandNumber() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(900000) + 100000)
}
func HashingPassword(password []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CompareHashedPassword(hashedPassword []byte, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return err
	}

	return nil
}
