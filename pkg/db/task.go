package db

import (
	"database/sql"
	"fmt"
	"time"
)
 
 type Task struct {
	ID 		string `json:"id"`
	Date	string `json:"date"`
	Title	string `json:"title"`
	Comment	string `json:"comment"`
	Repeat	string `json:"repeat"`
 }

 const dateFormat = "20060102"

 func Tasks (search string, limit int) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	switch {
		case search == "":
	
			rows, err = db.Query(
				`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`, 
				limit)

		case isDate(search):
			d, _ := time.Parse("02.01.2006", search)
			rows, err = db.Query(
				`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`,
				d.Format(dateFormat), limit)	
				
		default:
			like := "%" + search + "%"
			rows, err = db.Query(
				`SELECT id, date, title, comment, repeat FROM scheduler
				WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
				like, like, limit) 
	}

	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}

	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {

		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
 }

 func isDate(s string) bool {
	_, err := time.Parse("02.01.2006", s)
	return err == nil
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
		return 0, fmt.Errorf("get last insert id: %w", err)
	}
	return id, nil
 }

 func GetTask(id string) (*Task, error) {
	var t Task
	err := db.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler where id = ?`, id).
		Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	
	if err != nil {
		return nil, fmt.Errorf("tasks not found: %w", err)
	}

	return &t, nil
 }

 func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeate = ? WHERE id = ?`
	res ,err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	count, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}


	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
 }