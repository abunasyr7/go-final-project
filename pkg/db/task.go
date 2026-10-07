package db

import "fmt"
 
 type Task struct {
	ID 		string `json:"id"`
	Date	string `json:"date"`
	Title	string `json:"title"`
	Comment	string `json:"comment"`
	Repeat	string `json:"repeat"`
 }

 func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("add task: %w", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w, err")
	}
	return id, nil
 }