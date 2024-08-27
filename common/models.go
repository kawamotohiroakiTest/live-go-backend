package common

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var Todos = []Todo{
	{ID: 1, Title: "タスク 1", Completed: false},
	{ID: 2, Title: "2番目のタスク 2", Completed: true},
	{ID: 3, Title: "Sample Todo 3", Completed: false},
	{ID: 4, Title: "DB接続2", Completed: false},
}
