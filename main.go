package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var tasks []Task
var idCounter = 1

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	json.NewDecoder(r.Body).Decode(&newTask)
	newTask.ID = idCounter
	idCounter++
	newTask.Status = "pending"
	tasks = append(tasks, newTask)
	json.NewEncoder(w).Encode(newTask)
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, task := range tasks {
		if task.ID == id {
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(tasks)
}

// Handler to update an existing task
func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)                 // Get URL parameters
	id, err := strconv.Atoi(params["id"]) // Convert id from string to int
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Find the task by ID
	for i, task := range tasks {
		if task.ID == id {
			// Decode the request body to get updated task details
			var updatedTask Task
			err := json.NewDecoder(r.Body).Decode(&updatedTask)
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Update the task fields
			tasks[i].Title = updatedTask.Title
			tasks[i].Description = updatedTask.Description
			tasks[i].Status = updatedTask.Status

			// Return the updated task in response
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	// If task not found, return an error
	http.Error(w, "Task not found", http.StatusNotFound)
}

// Declare the main function
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks", getAllTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", getTaskByID).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	http.ListenAndServe(":8081", router)
}
