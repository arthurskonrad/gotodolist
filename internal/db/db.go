package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var (
	todos  []Todo
	nextID int
	mux    sync.Mutex
)

const filePath = "data/todos.json"

func Load() error {
	mux.Lock()
	defer mux.Unlock()

	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			todos = []Todo{}
			nextID = 1
			return nil
		}
		return err
	}

	err = json.Unmarshal(file, &todos)
	if err != nil {
		return err
	}

	// Achar o maior ID pra manter sequência
	max := 0
	for _, t := range todos {
		if t.ID > max {
			max = t.ID
		}
	}
	nextID = max + 1

	return nil
}

func Save() error {
	mux.Lock()
	defer mux.Unlock()

	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll("data", 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func All() []Todo {
	mux.Lock()
	defer mux.Unlock()

	// retorna uma cópia para não ter race condition
	c := make([]Todo, len(todos))
	copy(c, todos)
	return c
}

func Add(text string) {
	mux.Lock()
	defer mux.Unlock()

	todo := Todo{
		ID:   nextID,
		Text: text,
		Done: false,
	}

	todos = append(todos, todo)
	nextID++

	fmt.Println("✅ db.Add: tarefa adicionada:", todo)

	err := Save()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao salvar: %v\n", err)
	}
}

func Toggle(id int) {
	mux.Lock()
	defer mux.Unlock()

	for i, t := range todos {
		if t.ID == id {
			todos[i].Done = !t.Done
			break
		}
	}
	_ = Save()
}

func Delete(id int) {
	mux.Lock()
	defer mux.Unlock()

	newTodos := []Todo{}
	for _, t := range todos {
		if t.ID != id {
			newTodos = append(newTodos, t)
		}
	}
	todos = newTodos
	_ = Save()
}
