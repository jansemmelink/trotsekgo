package auth

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/jansemmelink/log"
	"github.com/jansemmelink/trotsek/lib/email"
	"github.com/jansemmelink/trotsek/lib/password"
)

//NewUser ...
func newUser(emailAddr string) (*user, error) {
	if err := email.CheckAddress(emailAddr); err != nil {
		return nil, log.Wrapf(err, "Invalid email address")
	}

	u := &user{
		id:    uuid.NewV1().String(),
		email: emailAddr,
		rpw:   "",
	}

	u.Reset()
	return u, nil
}

//IUser ...
type IUser interface {
	ID() string
	Email() string
	Reset() error
	Tpw() string
	SetPassword(tpw, rpw string) error
	Auth(rpw string) bool
}

type user struct {
	id     string
	email  string
	tpw    string
	tpwExp time.Time
	rpw    string
}

func (u user) ID() string {
	return u.id
}

func (u user) Email() string {
	return u.email
}

func (u user) Tpw() string {
	if time.Now().Before(u.tpwExp) {
		return u.tpw
	}
	return ""
}

func (u *user) Reset() error {
	//define temp password to reset with
	u.tpw = u.id[0:8]
	u.tpwExp = time.Now().Add(time.Hour)
	return nil
}

func (u *user) SetPassword(tpw, rpw string) error {
	if time.Now().After(u.tpwExp) {
		return log.Wrapf(nil, "temporary password expired")
	}
	if u.tpw != tpw {
		return log.Wrapf(nil, "temporary password mismatch")
	}
	if err := password.IsStrong(rpw); err != nil {
		return log.Wrapf(err, "password is too weak")
	}
	//accepted
	u.rpw = rpw
	u.tpw = ""
	u.tpwExp = time.Now()
	return nil
}

func (u user) Auth(rpw string) bool {
	if u.rpw != rpw {
		return false
	}
	return true
}
