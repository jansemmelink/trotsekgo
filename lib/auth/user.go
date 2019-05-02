package auth

//NewUser ...
func NewUser(name string) IUser {
	return user{
		name: name,
	}
}

//IUser ...
type IUser interface {
	Name() string
	Reset()
	Activate() error
	Login() (ISession, error)
}

type user struct {
	name string
}

func (u user) Name() string {
	return u.name
}
