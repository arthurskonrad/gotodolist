package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/arthurskonrad/gotodolist/internal/db"
)

type Todo struct {
	ID   int
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
	fmt.Println("✅ AddTodo foi chamado")

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("task")
	fmt.Println("➡️ Texto recebido:", text)

	if text == "" {
		http.Error(w, "Texto da tarefa vazio", http.StatusBadRequest)
		return
	}

	db.Add(text)
	fmt.Println("📦 Tarefas após add:", db.All())

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "layout.html"), // obrigatório!
		filepath.Join("internal", "templates", "index.html"),
	)
	if err != nil {
		http.Error(w, "Erro no template: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("❌ Erro no ParseFiles:", err)
		return
	}

	fmt.Println("👀 Dados enviados ao template:", splitTodos())

	err = tmpl.ExecuteTemplate(w, "content", splitTodos())
	if err != nil {
		http.Error(w, "Erro ao renderizar template: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("❌ Erro ao renderizar template:", err)
		return
	}

	fmt.Println("✅ Template renderizado com sucesso")
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	db.Delete(id)

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "layout.html"),
		filepath.Join("internal", "templates", "index.html"),
	)
	if err != nil {
		http.Error(w, "Erro ao carregar templates", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "content", splitTodos())
	if err != nil {
		http.Error(w, "Erro ao renderizar template: "+err.Error(), http.StatusInternalServerError)
	}
}

func ToggleDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	db.Toggle(id)

	tmpl, err := template.ParseFiles(
		filepath.Join("internal", "templates", "layout.html"),
		filepath.Join("internal", "templates", "index.html"),
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

func splitTodos() TodoViewData {
	all := db.All()

	var pending []db.Todo
	var completed []db.Todo

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
