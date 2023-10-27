package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const MinSecretkeySize = 32

type JwtMaker struct {
	secretKey string

}

func NewJwtMaker(secretkey string) (Maker, error)  {
	if len(secretkey) < MinSecretkeySize {
		return nil, fmt.Errorf("Invalid key size: must be atleast %d character.", MinSecretkeySize)
	}
	return &JwtMaker{secretKey: secretkey},nil
}

func ( j *JwtMaker)CreateToken(username string, duration time.Duration) (string,error){
	payload, err := NewPayload(username,duration)
	if err != nil{
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,payload)
	return token.SignedString([]byte(j.secretKey))
}
func ( j *JwtMaker)VerifyToken(token string) (*Payload, error){
	keyFunc := func(token *jwt.Token) (interface{}, error){
	_,ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil,ErrInvalidToken
	}
		return []byte(j.secretKey),nil
	}
	jwtToken, err := jwt.ParseWithClaims(token,&Payload{},keyFunc)
	if err != nil {
		if verr,ok := err.(*jwt.ValidationError);ok{
			if errors.Is(verr.Inner, ErrExpiredToken){
				return nil,ErrExpiredToken
			}
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload);
	if !ok{
		return nil, ErrInvalidToken
	}
	return payload,nil
}