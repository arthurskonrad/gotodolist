package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var (
	todos  []Todo
	nextID string
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
			nextID = uuid.New().String()
			return nil
		}
		return err
	}

	err = json.Unmarshal(file, &todos)
	if err != nil {
		return err
	}

	nextID = uuid.New().String()

	return nil
}

func Save() error {
	mux.Lock()
	defer mux.Unlock()

	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		fmt.Println("db.Save: erro ao indentar o json:", err)
		return err
	}

	err = os.MkdirAll("data", 0755)
	if err != nil {
		fmt.Println("db.Save: erro ao criar pasta:", err)
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		fmt.Println("db.Save: erro ao escrever no arquivo:", err)
		return err
	}

	return nil
}

func All() []Todo {
	mux.Lock()
	defer mux.Unlock()

	if todos == nil {
		return []Todo{}
	}

	c := make([]Todo, len(todos))
	copy(c, todos)
	return c
}

func Add(text string) Todo {
	todo := Todo{
		ID:   nextID,
		Text: text,
		Done: false,
	}

	mux.Lock()
	todos = append(todos, todo)
	nextID = uuid.New().String()
	mux.Unlock()

	_ = Save()

	return todo
}

func Toggle(id string) {
	mux.Lock()
	defer mux.Unlock()

	for i, t := range todos {
		if t.ID == id {
			todos[i].Done = !t.Done
			break
		}
	}
}

func Delete(id string) {
	mux.Lock()
	defer mux.Unlock()

	filtered := make([]Todo, 0, len(todos))
	for _, t := range todos {
		if t.ID != id {
			filtered = append(filtered, t)
		}
	}
	todos = filtered
	nextID = uuid.New().String()
}
