package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TodoPayLoad struct {
	Desc string `json:"description"`
}

var TodosStore = TodoStore{
	todos: []Todo{{1, "Read a book", false},
		{2, "Exercise", true}},
}

func main() {

	fmt.Println("Building REST APIs in Go 1.22!")

	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/api/todos", getTodoHandler).Methods("GET")
	router.HandleFunc("/api/todos/{id}", getTodoDetailHandler).Methods("GET")
	router.HandleFunc("/api/todos", addTodoHandler).Methods("POST")
	router.HandleFunc("/api/todos/{id}/mark_completed", markCompleteHandler).Methods("PUT")
	router.HandleFunc("/api/todos/{id}/mark_incomplete", markIncompleteHandler).Methods("PUT")
	router.HandleFunc("/api/todos/{id}", updateTodoHandler).Methods("PUT")
	router.HandleFunc("/api/todos/{id}", deleteTodoHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe("localhost:9000", router))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	writeMessage(w, http.StatusOK, "Hello Venkat1")
}

func getTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("completed") == "true" {
		writeJson(w, http.StatusOK, TodosStore.GetByStatus(true))
	} else if r.URL.Query().Get("completed") == "false" {
		writeJson(w, http.StatusOK, TodosStore.GetByStatus(false))
	} else {
		writeJson(w, http.StatusOK, TodosStore.GetAll())
	}
}

func getTodoDetailHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"]) //string to intger conversion
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid id")
		return
	}
	todo, err1 := TodosStore.GetTodoDetail(id)
	if err1 != nil {
		writeMessage(w, http.StatusBadRequest, err1.Error())
		return
	}
	writeJson(w, http.StatusOK, todo)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	var payLoad TodoPayLoad
	json.NewDecoder(r.Body).Decode(&payLoad)
	TodosStore.AddTodo(payLoad.Desc)
	writeMessage(w, http.StatusOK, "Successfully addded")
}

func markCompleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"]) //string to intger conversion
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid id")
		return
	}

	error := TodosStore.SetCompleted(id, true)
	if error != nil {
		writeMessage(w, http.StatusNotFound, error.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully marked completed")
}

func markIncompleteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"]) //string to intger conversion
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid id")
		return
	}

	error := TodosStore.SetCompleted(id, false)
	if error != nil {
		writeMessage(w, http.StatusNotFound, error.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully marked incomplete")
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"]) //string to intger conversion
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid id")
		return
	}
	var payLoad TodoPayLoad
	json.NewDecoder(r.Body).Decode(&payLoad)

	error := TodosStore.Update(id, payLoad.Desc)
	if error != nil {
		writeMessage(w, http.StatusNotFound, error.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully updated")
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"]) //string to intger conversion
	if err != nil {
		writeMessage(w, http.StatusBadRequest, "Invalid id")
		return
	}

	error := TodosStore.Delete(id)
	if error != nil {
		writeMessage(w, http.StatusNotFound, error.Error())
		return
	}
	writeMessage(w, http.StatusOK, "Successfully deleted")
}

func writeJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeMessage(w http.ResponseWriter, status int, message string) {
	writeJson(w, status, map[string]string{
		"message": message,
	})
}
