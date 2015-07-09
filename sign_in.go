package kolkata

import (
	"github.com/likestripes/pacific"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type SignIn struct {
	Token        string          `gorm:"column:signin_id"`
	PersonId     int64
	PasswordHash []byte
	Unsafe       string          `datastore:"-" sql:"-" json:"-"`
	Context      pacific.Context `datastore:"-" sql:"-" json:"-"`
}

func (signin SignIn) Auth() (int64, error) {
	err := bcrypt.CompareHashAndPassword(signin.PasswordHash, []byte(signin.Unsafe))
	if err == nil {
		return signin.PersonId, nil
	}
	return 0, err
}

func (signin *SignIn) Get() error {
	query := pacific.Query{
		Context:   signin.Context,
		Kind:      "SignIn",
		KeyString: signin.Token,
	}
	return query.Get(signin)
}

func (signin SignIn) Save() (SignIn, error) {

	signin.Token = strings.ToLower(signin.Token)
	signin.PasswordHash = hashString(signin.Unsafe)

	query := pacific.Query{
		Context:   signin.Context,
		Kind:      "SignIn",
		KeyString: signin.Token,
	}

	err := query.Put(&signin)

	return signin, err
}
