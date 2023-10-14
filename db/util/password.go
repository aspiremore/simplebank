package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error)  {
	hash, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil{
		return "",fmt.Errorf("Error converting input to hashed password %w", err)
	}
	return string(hash), nil
}

func CheckPassword(hashedpassword string, userpassword string) error  {
	return bcrypt.CompareHashAndPassword([]byte(hashedpassword),[]byte(userpassword))
}
