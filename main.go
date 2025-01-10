package main

import (
	"html/template"
	"net/http"
)

type Todo struct {
	Title string
	Done bool
}

type TodoPageData struct {
	PageTitle string
	Todos []Todo
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/index.html"))

	data := TodoPageData{
		PageTitle: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl.Execute(w, data)
}

func main () {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

// git add . git commit -m "msg" git push