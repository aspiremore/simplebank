package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoMaker struct {
	paseto *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize{
		return nil, fmt.Errorf("Invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	return &PasetoMaker{
		paseto: paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	},nil
	
}

func ( p *PasetoMaker)CreateToken(username string, duration time.Duration) (string,error){
	payload, err := NewPayload(username,duration)
	if err != nil{
		return "", err
	}
	return p.paseto.Encrypt(p.symmetricKey,payload,nil)
}
func ( p *PasetoMaker)VerifyToken(token string) (*Payload, error){
	payload := &Payload{}
	err := p.paseto.Decrypt(token,p.symmetricKey,payload,nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil,err
	}
	return payload, nil
}
