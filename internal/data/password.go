package data

import "golang.org/x/crypto/bcrypt"

type Password struct {
	Plaintext *string
	hash      []byte
}

func (p *Password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.hash = hash
	p.Plaintext = &plaintext

	return nil
}
