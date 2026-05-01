package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
	Genre string `json:"genre"`
}

var db *sql.DB

func connectDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("БД недоступна: ", err)
	}

	fmt.Println("Подключение к PostgreSQL успешно!")
}

func main() {
	connectDB()

	http.HandleFunc("/", handle)
	http.HandleFunc("/book", handleBook)
	http.HandleFunc("/book/", handleBookID)

	fmt.Println("🚀 Сервер запущен на :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Сервер запущен")
}

func handleBook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		rows, err := db.Query("SELECT id, title, year, genre FROM books")
		if err != nil {
			http.Error(w, "Ошибка получения книг", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var books []book
		for rows.Next() {
			var b book
			err := rows.Scan(&b.ID, &b.Title, &b.Year, &b.Genre)
			if err != nil {
				http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
				return
			}
			books = append(books, b)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)

	case http.MethodPost:
		var b book
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			http.Error(w, "Ошибка чтения JSON", http.StatusBadRequest)
			return
		}

		err = db.QueryRow(
			"INSERT INTO books (title, year, genre) VALUES ($1, $2, $3) RETURNING id",
			b.Title, b.Year, b.Genre,
		).Scan(&b.ID)
		if err != nil {
			http.Error(w, "Ошибка добавления книги", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(b)
	}
}

func handleBookID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodPut:
		var b book
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			http.Error(w, "Ошибка чтения JSON", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"UPDATE books SET title=$1, year=$2, genre=$3 WHERE id=$4",
			b.Title, b.Year, b.Genre, id,
		)
		if err != nil {
			http.Error(w, "Ошибка обновления", http.StatusInternalServerError)
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(w, "Книга не найдена", http.StatusNotFound)
			return
		}

		b.ID = id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(b)

	case http.MethodDelete:
		result, err := db.Exec("DELETE FROM books WHERE id=$1", id)
		if err != nil {
			http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(w, "Книга не найдена", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
