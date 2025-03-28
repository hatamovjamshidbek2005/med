package password_hash

import "golang.org/x/crypto/bcrypt"

func CompareHashedPassword(hashedPassword []byte, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return err
	}

	return nil
}

func HashingPassword(password []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
