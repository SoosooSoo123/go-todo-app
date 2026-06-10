package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	InitDB()
	http.HandleFunc("/tasks", handleTasks)
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		getTasks(w, r)
	case "POST":
		createTask(w, r)
	case "PUT":
		updateTask(w, r)
	case "DELETE":
		deleteTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	if result := DB.Find(&todos); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todos)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if result := DB.Create(&todo); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "id نامعتبر", http.StatusBadRequest)
		return
	}
	var todo Todo
	if result := DB.First(&todo, id); result.Error != nil {
		http.Error(w, "تسک پیدا نشد", http.StatusNotFound)
		return
	}
	todo.Completed = !todo.Completed
	DB.Save(&todo)
	json.NewEncoder(w).Encode(todo)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "id نامعتبر", http.StatusBadRequest)
		return
	}
	if result := DB.Delete(&Todo{}, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}