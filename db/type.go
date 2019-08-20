package db

type Todo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
	Done bool   `json:"done"`
}

type Welcom struct {
	Name    string `json:"name"`
	Version int64  `json:"version"`
}
