package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/arthurskonrad/gotodolist/internal/db"
)

type Todo struct {
	ID   string
	Text string
	Done bool
}

type TodoViewData struct {
	Pending   []db.Todo
	Completed []db.Todo
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "layout.html"),
		filepath.Join("internal", "templates", "index.html"),
		filepath.Join("internal", "templates", "item.html"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", splitTodos())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("task")
	if text == "" {
		http.Error(w, "Texto inválido", http.StatusBadRequest)
		return
	}

	newTodo := db.Add(text)

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "item.html"),
	)
	if err != nil {
		http.Error(w, "Erro ao carregar template do item", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "item", newTodo)
	if err != nil {
		http.Error(w, "Erro ao renderizar item", http.StatusInternalServerError)
	}
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	db.Delete(idStr)

	if err := db.Save(); err != nil {
		http.Error(w, "Erro ao salvar dados", http.StatusInternalServerError)
		return
	}
}

func ToggleDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	db.Toggle(idStr)

	if err := db.Save(); err != nil {
		http.Error(w, "Erro ao salvar dados", http.StatusInternalServerError)
		return
	}

	renderTodos(w)
}

func splitTodos() TodoViewData {
	all := db.All()

	pending := make([]db.Todo, 0)
	completed := make([]db.Todo, 0)

	for _, t := range all {
		if t.Done {
			completed = append(completed, t)
		} else {
			pending = append(pending, t)
		}
	}

	return TodoViewData{
		Pending:   pending,
		Completed: completed,
	}
}

func renderTodos(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "layout.html"),
		filepath.Join("internal", "templates", "index.html"),
		filepath.Join("internal", "templates", "item.html"),
	)
	if err != nil {
		http.Error(w, "Erro ao carregar templates", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "content", splitTodos())
	if err != nil {
		http.Error(w, "Erro ao renderizar template", http.StatusInternalServerError)
	}
}
