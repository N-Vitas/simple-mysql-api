package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type Con struct {
	db         *sql.DB
	accessLock *sync.RWMutex
	host       string
	port       int
	user       string
	pass       string
	database   string
}

func Init(host string, port int, user string, pass string, database string) *Con {
	return &Con{
		db:         nil,
		accessLock: &sync.RWMutex{},
		host:       host,
		port:       port,
		user:       user,
		pass:       pass,
		database:   database,
	}
}

func (s *Con) GetDb() *sql.DB {
	s.accessLock.RLock()
	existing := s.db
	s.accessLock.RUnlock()
	var err error
	// Clone and return if sessions exists
	if existing != nil {
		err = existing.Ping()
		if err != nil {
			existing = nil
		} else {
			//info("return connect to %s : %s",alias,sessionId)
			return existing
		}
	}
	// Get timeout from configuration
	s.accessLock.Lock()
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", s.user, s.pass, s.host, s.port, s.database)
	s.db, err = sql.Open("mysql", connString)
	if err != nil {
		panic(err.Error())
	}
	s.accessLock.Unlock()
	// s.sessions = newSession
	return s.db
}
