package models

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var todos = []Todo{}
var lastID = 0

func GetTodos() []Todo {
	return todos
}

func AddTodo(title string) Todo {
	lastID++
	todo := Todo{
		ID:    lastID,
		Title: title,
		Done:  false,
	}
	todos = append(todos, todo)
	return todo
}

func ToggleTodo(id int) (Todo, bool) {
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Done = !todos[i].Done
			return todos[i], true
		}
	}
	return Todo{}, false
}

func DeleteTodo(id int) bool {
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return true
		}
	}
	return false
}
