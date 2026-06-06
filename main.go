package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/glebarez/go-sqlite" // درایور دیتابیس SQLite
)

// متغیر سراسری برای نگه داشتن اتصال دیتابیس
var db *sql.DB

func main() {
	var err error
	// ۱. اتصال به دیتابیس (اگر فایل وجود نداشته باشد، آن را می‌سازد)
	db, err = sql.Open("sqlite", "todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ۲. ساخت جدول وظایف (اگر از قبل ساخته نشده باشد)
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		done INTEGER DEFAULT 0
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to your Go To-Do List (Database Version)!")

	for {
		fmt.Println("\n--- MENU ---")
		fmt.Println("1. View Tasks")
		fmt.Println("2. Add Task")
		fmt.Println("3. Mark Task as Done")
		fmt.Println("4. Delete Task")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option (1-5): ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			listTasks()
		case "2":
			fmt.Print("Enter task description: ")
			scanner.Scan()
			taskText := strings.TrimSpace(scanner.Text())
			if taskText != "" {
				addTask(taskText)
			} else {
				fmt.Println("Task cannot be empty!")
			}
		case "3":
			listTasks()
			fmt.Print("Enter the ID of the task to mark as done: ")
			scanner.Scan()
			idStr := strings.TrimSpace(scanner.Text())
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID. Please enter a number.")
				continue
			}
			markDone(id)
		case "4":
			listTasks()
			fmt.Print("Enter the ID of the task to delete: ")
			scanner.Scan()
			idStr := strings.TrimSpace(scanner.Text())
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID. Please enter a number.")
				continue
			}
			deleteTask(id)
		case "5":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please select 1-5.")
		}
	}
}

// افزودن تسک به دیتابیس (INSERT)
func addTask(taskText string) {
	query := `INSERT INTO tasks (task, done) VALUES (?, 0)`
	_, err := db.Exec(query, taskText)
	if err != nil {
		fmt.Println("Error adding task:", err)
		return
	}
	fmt.Printf("Added: \"%s\"\n", taskText)
}

// نمایش تسک‌ها از دیتابیس (SELECT)
func listTasks() {
	query := `SELECT id, task, done FROM tasks`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		return
	}
	defer rows.Close()

	hasTasks := false
	fmt.Println("\n--- YOUR TASKS ---")
	for rows.Next() {
		hasTasks = true
		var id int
		var task string
		var done int
		err = rows.Scan(&id, &task, &done)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		status := "[ ]"
		if done == 1 {
			status = "[✓]"
		}
		fmt.Printf("%d. %s %s\n", id, status, task)
	}

	if !hasTasks {
		fmt.Println("Your to-do list is empty!")
	}
}

// بروزرسانی وضعیت تسک در دیتابیس (UPDATE)
func markDone(id int) {
	query := `UPDATE tasks SET done = 1 WHERE id = ?`
	result, err := db.Exec(query, id)
	if err != nil {
		fmt.Println("Error updating task:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Printf("Task with ID %d not found.\n", id)
	} else {
		fmt.Printf("Task %d marked as done!\n", id)
	}
}

// حذف تسک از دیتابیس (DELETE)
func deleteTask(id int) {
	query := `DELETE FROM tasks WHERE id = ?`
	result, err := db.Exec(query, id)
	if err != nil {
		fmt.Println("Error deleting task:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Printf("Task with ID %d not found.\n", id)
	} else {
		fmt.Printf("Task %d deleted successfully.\n", id)
	}
}
