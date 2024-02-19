package data

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

const cookie_name string = "jb-cookie"

const inputDate = "2006-01-02 15:04:05 -0700 MST"

var sessionManager *SessionManager

func init() {
	sessionManager = &SessionManager{sessions: make(map[string]Session)}
	go RunSessionKiller()
}

type SessionManager struct {
	sync.RWMutex
	sessions map[string]Session
}

func (s *SessionManager) GetSession(key string) (Session, bool) {
	s.RLock()
	f, ok := s.sessions[key]
	s.RUnlock()
	return f, ok
}

func (s *SessionManager) SetSession(key string, obj Session) {
	s.Lock()
	s.sessions[key] = obj
	s.Unlock()
}

func (s *SessionManager) CreateSession(key string, user User) {
	session := Session{}
	session.Start = time.Now()
	session.User = user

	s.SetSession(key, session)
}

func (s *SessionManager) UpdateSession(r *http.Request, user User) error {

	c, err := r.Cookie(cookie_name)
	if err != nil {
		return err
	}

	session, ok := s.GetSession(c.Value)
	if !ok {
		return errors.New("no session found")
	}
	session.Start = time.Now()
	session.User = user

	s.SetSession(c.Value, session)
	return nil
}

func (s *SessionManager) GetUser(key string) (User, error) {
	sess, ok := s.GetSession(key)
	if !ok {
		return User{}, errors.New("no session found")
	}

	return sess.User, nil
}

func (s *SessionManager) DeleteSession(key string) {
	s.Lock()
	delete(s.sessions, key)
	s.Unlock()
}

type Session struct {
	Start time.Time
	User  User
}

func RunSessionKiller() {
	for {
		sessionManager.Lock()
		now := time.Now()
		for k, v := range sessionManager.sessions {
			if now.Sub(v.Start) > 10*time.Minute {
				delete(sessionManager.sessions, k)
			}
		}
		sessionManager.Unlock()
		time.Sleep(time.Minute) //important unlock before sleeping
	}
}
