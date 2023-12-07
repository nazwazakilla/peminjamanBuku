package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Book represents a book entity.
type Book struct {
	Title     string
	Author    string
	Available bool
	Stock     int
}

var bookList = []Book{
	{"Si Kancil", "Tira Ikranegara", true, 5},
	{"Tentang Kamu", "Tere Liye", true, 3},
	{"7 Habits of Highly Effective People", "Stephen Covey", false, 0},
	{"The Art Of Thinking Clearly", "Rolf Dobelli", true, 1},
	{"Pemrograman Web Dasar", "Rintho Rare Rerung", false, 0},
	{"Dikta dan Hukum", "Dhia'an Farah", true, 12},
	{"Nihonggo Mantappu", "Jerome Polin", false, 0},
	{"Kambing Jantan", "Raditya Dika", true, 8},
}

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, bookList)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	results := []Book{}

	for _, book := range bookList {
		if contains(book.Title, query) || contains(book.Author, query) {
			results = append(results, book)
		}
	}

	tpl.Execute(w, results)
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func borrowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	book, found := findBookByTitle(title)
	if !found {
		http.NotFound(w, r)
		return
	}

	if book.Stock <= 0 {
		fmt.Fprintf(w, `<script>alert('This book is not available or out of stock.'); window.location.href='/';</script>`)
		return
	}

	renderTemplate(w, "borrow.html", book)
}

func confirmBorrowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	book, found := findBookByTitle(title)
	if !found {
		http.NotFound(w, r)
		return
	}

	if !book.Available || book.Stock <= 0 {
		// Display an alert on the client side
		fmt.Fprintf(w, `<script>alert('This book is not available or out of stock.'); window.location.href='/';</script>`)
		return
	}

	if r.Method == http.MethodPost {
		// Handle the form submission, e.g., store the data in a database.
		// For simplicity, we just display a thank you message in this example.
		durationStr := r.FormValue("duration")
		duration, err := strconv.Atoi(durationStr)
		if err != nil || duration <= 0 {
			// Durasi tidak valid, kembalikan ke halaman borrow.html
			http.Redirect(w, r, "/borrow/"+title, http.StatusSeeOther)
			return
		}

		renderTemplate(w, "thankyou.html", nil)
		book.Stock-- // Kurangi stok
		if book.Stock == 0 {
			book.Available = false
		}
	}
}

func findBookByTitle(title string) (*Book, bool) {
	for i := range bookList {
		if bookList[i].Title == title {
			return &bookList[i], true
		}
	}
	return nil, false
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()

	// Handle the main index page
	r.HandleFunc("/", indexHandler).Methods("GET")

	// Handle the search route
	r.HandleFunc("/search", searchHandler).Methods("GET")

	http.Handle("/", r)

	r.HandleFunc("/borrow/{title}", borrowHandler).Methods("GET")
	r.HandleFunc("/confirm-borrow/{title}", confirmBorrowHandler).Methods("POST")

	fmt.Println("Server is running on :9000...")
	http.ListenAndServe(":9000", nil)
}
