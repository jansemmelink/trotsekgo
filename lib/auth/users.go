package auth

import (
	"sync"

	"github.com/jansemmelink/log"
)

//NewUsers ...
func NewUsers() IUsers {
	udb := &users{
		emailList: make(map[string]*user),
		idList:    make(map[string]*user),
		nextID:    1,
	}
	return udb
}

//IUsers is a db of users
type IUsers interface {
	Add(email string) (IUser, error)
	GetByEmail(email string) IUser
	GetByID(id string) IUser
	//Delete(IUser)
}

//users implement IUsers
type users struct {
	mutex     sync.Mutex
	emailList map[string]*user
	idList    map[string]*user
	nextID    int
}

//Add ...
func (udb *users) Add(email string) (IUser, error) {
	if udb.GetByEmail(email) != nil {
		return nil, log.Wrapf(nil, "User with email %s already registered.", email)
	}

	u, err := newUser(email)
	if err != nil {
		return nil, err
	}

	//created, add:
	udb.emailList[u.email] = u
	udb.idList[u.id] = u
	return u, nil
}

func (udb *users) GetByEmail(emailAddr string) IUser {
	if u, ok := udb.emailList[emailAddr]; ok {
		return u
	}
	return nil
}

func (udb *users) GetByID(id string) IUser {
	if u, ok := udb.idList[id]; ok {
		return u
	}
	return nil
}
