package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	user     = "postgres"
	password = "12345S"
	host     = "localhost"
	port     = 5432
	dbname   = "Sandugash"
)

func connectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

type Task struct {
	ID        int
	Name      string
	Completed bool
}

func main() {
	connectionString()

	db, err := sql.Open("postgres", connectionString())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = createTask("sgse")
	if err != nil {
		panic(err)
	}
}
func createTask(name string) error {
	db, err := sql.Open("postgres", connectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if task with same name already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE name = $1", name).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("task with name '%s' already exists", name)
	}

	// Insert new task
	_, err = db.Exec("INSERT INTO tasks (name) VALUES ($1)", name)
	if err != nil {
		return err
	}

	fmt.Println("Task created successfully")
	return nil
}
func readTasks() ([]Task, error) {
	db, err := sql.Open("postgres", connectionString())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func updateTaskCompleted(id int, completed bool) error {
	db, err := sql.Open("postgres", connectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE tasks SET completed = $1 WHERE id = $2", completed, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Println("Task updated successfully")
	return nil
}
func deleteTask(id int) error {
	db, err := sql.Open("postgres", connectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Println("Task deleted successfully")
	return nil
}
