package auth

import (
	"regexp"

	"github.com/jansemmelink/log"
	"github.com/jansemmelink/trotsek/lib/router"
)

//New ...
func New() IAuth {
	a := &auth{
		user: make(map[string]IUser),
	}
	return a
}

//IAuth implements authetication operations
type IAuth interface {
	RegisterOper() router.IOper
	Register(uname string) (IUser, error)
	ResetOper() router.IOper
	ActivateOper() router.IOper
	LoginOper() router.IOper
	Login(uname string) (IUser, error)
	LogoutOper() router.IOper
}

type auth struct {
	user map[string]IUser
}

func (a *auth) RegisterOper() router.IOper { return Register{a: a} }
func (a *auth) ResetOper() router.IOper    { return Reset{a: a} }
func (a *auth) ActivateOper() router.IOper { return Activate{a: a} }
func (a *auth) LoginOper() router.IOper    { return Login{a: a} }
func (a *auth) LogoutOper() router.IOper   { return Logout{a: a} }

func (a *auth) Register(uname string) (IUser, error) {
	if _, ok := a.user[uname]; ok {
		return nil, log.Wrapf(nil, "user %s already exists", uname)
	}
	a.user[uname] = NewUser(uname)
	return a.user[uname], nil
}

func (a *auth) Login(uname string) (IUser, error) {
	return nil, log.Wrapf(nil, "NYI")
}

//Register implements router.IOper
type Register struct {
	a     IAuth
	Uname string `json:"uname"`
}

//New ...
func (r Register) New() router.IOper {
	return &Register{
		a: r.a,
	}
}

//Validate ...
func (r Register) Validate() error {
	if len(r.Uname) < 1 {
		return log.Wrapf(nil, "Missing uname")
	}
	if !validUname.MatchString(r.Uname) {
		return log.Wrapf(nil, "Invalid uname=\"%s\"", r.Uname)
	}
	return nil
}

//Exec ...
func (r Register) Exec() (interface{}, error) {
	u, err := r.a.Register(r.Uname)
	if err != nil {
		return nil, err
	}
	return u, nil
}

//Reset ...
type Reset struct{ a IAuth }

//New ...
func (o Reset) New() router.IOper {
	return &Reset{}
}

//Validate ...
func (o Reset) Validate() error {
	return nil
}

//Exec ...
func (o Reset) Exec() (interface{}, error) {
	return nil, nil
}

//Activate ...
type Activate struct{ a IAuth }

//New ...
func (o Activate) New() router.IOper {
	return &Activate{}
}

//Validate ...
func (o Activate) Validate() error {
	return nil
}

//Exec ...
func (o Activate) Exec() (interface{}, error) {
	return nil, nil
}

//Logout ...
type Logout struct{ a IAuth }

//New ...
func (o Logout) New() router.IOper {
	return &Logout{}
}

//Validate ...
func (o Logout) Validate() error {
	return nil
}

//Exec ...
func (o Logout) Exec() (interface{}, error) {
	return nil, nil
}

//Login implemented router.IOper
type Login struct {
	a IAuth
}

//New ...
func (l Login) New() router.IOper {
	return Login{}
}

//Validate ...
func (l Login) Validate() error {
	return nil
}

//Exec ...
func (l Login) Exec() (interface{}, error) {
	return nil, log.Wrapf(nil, "NYI")
}

var validUname *regexp.Regexp

func init() {
	validUname = regexp.MustCompile(`^[a-z][a-z0-9.@_-]*`)
}
