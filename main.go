package main

import (
	"flag"
	"log"
	"net/http"
	"todo-api/db"
)

func main() {
	host := flag.String("host", "127.0.0.1", "a string")
	port := flag.Int("port", 3306, "an int")
	user := flag.String("user", "root", "a string")
	pass := flag.String("password", "123", "a string")
	database := flag.String("database", "todo", "a string")
	flag.Parse()

	app := db.Init(*host, *port, *user, *pass, *database)
	app.GetDb()
	http.HandleFunc("/", app.Default)

	log.Println("Go!")
	http.ListenAndServe(":8080", nil)
}
