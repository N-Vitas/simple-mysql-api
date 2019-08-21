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
	pass := flag.String("password", "", "a string")
	database := flag.String("database", "todo", "a string")
	flag.Parse()

	app := db.Init(*host, *port, *user, *pass, *database)
	http.HandleFunc("/", app.Default)
	http.HandleFunc("/todo", app.Todos)
	// for i := 1; i < 50; i++ {
	// 	t := &db.Todo{
	// 		Name: fmt.Sprintf("Задача номер %d", i),
	// 		Date: time.Now().String(),
	// 		Done: false,
	// 	}
	// 	app.SaveTodo(t)
	// }
	log.Println("Go!")
	http.ListenAndServe(":3000", nil)
}
