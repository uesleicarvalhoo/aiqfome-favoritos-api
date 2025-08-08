package ioc

import (
	"sync"

	"github.com/uesleicarvalhoo/aiqfome/config"
	"github.com/uesleicarvalhoo/aiqfome/pkg/password"
)

var (
	passwordHasher     password.Hasher
	passwordHasherOnce sync.Once
)

func PasswordHasher() password.Hasher {
	passwordHasherOnce.Do(func() {
		cost := config.GetInt("PASSWORD_HASHSER_CRYPT_COST")
		passwordHasher = password.NewBcryptHasher(cost)
	})

	return passwordHasher
}
