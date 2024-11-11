package token

import (
	"github.com/o1egl/paseto"
	"time"
)

type PasetoToken struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func (p PasetoToken) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := p.paseto.Encrypt(p.symmetricKey, payload, nil)
	if err != nil {
		return "", payload, err
	}
	return token, payload, nil
}

func (p PasetoToken) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func NewPasetoToken() (*PasetoToken, error) {
	//if len(symmetricKey) != chacha20poly1305.KeySize {
	//	return nil, fmt.Errorf("symmetric key length must be 20 bytes")
	//}

	maker := &PasetoToken{
		symmetricKey: []byte("12345678912345678912345678923456"),
		paseto:       paseto.NewV2(),
	}
	return maker, nil
}
