package auth

import (
	"time"
)

var (
	sessions = make(map[string]ISession)
	sessionsMutex = sync.Mutex{}
)

const expireDuration = time.Minute * 5

//NewSession ...
func NewSession(u IUser) ISession {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	s := &session{
		id: uuid.NewV1(),
		user:   u,
		start:  time.Now(),
		expire: time.Now().Add(expireDuration),
	}
	sessions[s.id] = s
	return s
}

//ISession ...
type ISession interface {
	ID() string
	User() IUser
	Logout()
}

type session struct {
	user   IUser
	start  time.Time
	expire time.Time
}
