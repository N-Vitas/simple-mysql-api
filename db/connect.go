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
	return s.db
}

func (s *Con) GetTodos(start int64, finish int64) []Todo {
	todos := []Todo{}
	rows, err := s.GetDb().Query("SELECT id, name, date, done FROM list ORDER BY ID ASC LIMIT ?, ?", start, finish)
	if err != nil {
		s.Error("%v", err)
		return todos
	}
	for rows.Next() {
		t := Todo{}
		err = rows.Scan(&t.Id, &t.Name, &t.Date, &t.Done)
		if err != nil {
			continue
		}
		todos = append(todos, t)
	}
	return todos
}

func (s *Con) GetTodo(id int64) (Todo, bool) {
	t := Todo{}
	err := s.GetDb().QueryRow("SELECT id, name, date, done FROM list WHERE id=?", id).Scan(&t.Id, &t.Name, &t.Date, &t.Done)
	if err != nil {
		s.Error("QueryRow %v", err)
		return t, false
	}
	return t, true
}
func (s *Con) ChangeTodo(t Todo) bool {
	_, err := s.GetDb().Exec("UPDATE `list` SET `name`=?,`date`=?,`done`=? WHERE id=?", t.Name, t.Date, t.Done, t.Id)
	if err != nil {
		s.Error("QueryRow %v", err)
		return false
	}
	return true
}
func (s *Con) DelTodo(t Todo) bool {
	_, err := s.GetDb().Exec("DELETE FROM `list` WHERE id=?", t.Id)
	if err != nil {
		s.Error("QueryRow %v", err)
		return false
	}
	return true
}
func (s *Con) SaveTodo(todo *Todo) bool {
	r, err := s.GetDb().Exec("INSERT INTO `list`(`name`, `date`, `done`) VALUES (?,?,?)", todo.Name, todo.Date, todo.Done)
	if err != nil {
		s.Error("%s", err.Error())
		return false
	}
	todo.Id, _ = r.LastInsertId()
	return true
}

func (s *Con) Error(template string, arg ...interface{}) {
	fmt.Printf(template+"\n", arg...)
}
