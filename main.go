package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
	Genre string `json:"genre"`
}

var books = []book{}
var newBook book

func main() {

	http.HandleFunc("/", handle)
	http.HandleFunc("/book", handleBook)
	http.HandleFunc("/book/", handleBookId)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err, "Ошибка сервера")
	}

}
func handleBook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(books)
		if err != nil {
			http.Error(w, "Ошибка при отправке", http.StatusInternalServerError)
		}

	case http.MethodPost:
		item, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(item, &newBook)
		if err != nil {
			http.Error(w, "ошибка обработки данных в JSON", http.StatusBadRequest)
		}
		newBook.ID = len(books) + 1
		fmt.Println("Raw body:", string(item))
		books = append(books, newBook)
		fmt.Fprintf(w, "Книга добавлена:%v\n", books)
	}

}
func handleBookId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:

		var updateBook book

		item, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения", http.StatusBadRequest)
			return
		}
		data := r.URL.Path
		parts := strings.Split(data, "/")
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "Ошибка чтения ID", http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(item, &updateBook)
		if err != nil {
			http.Error(w, "Ошибка JSON", http.StatusBadRequest)
			return
		}
		updated := false
		for i, b := range books {
			if b.ID == id {
				updated = true
				updateBook.ID = id
				books[i] = updateBook
			}
		}

		if !updated {
			http.Error(w, "Книга с таким ID нету", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Книга добавлена:%v\n", books)

	case http.MethodDelete:
		data := r.URL.Path
		parts := strings.Split(data, "/")
		id, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "Ошибка ссылки ID", http.StatusBadRequest)
			return
		}
		updated := false
		for i, b := range books {
			if b.ID == id {
				books = append(books[:i], books[i+1:]...)
				updated = true

			}

		}
		if !updated {
			http.Error(w, "Книга с таким ID нету", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}

func handle(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Сервер запущен")

}
