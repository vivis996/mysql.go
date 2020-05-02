package main

import (
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var tmpl = template.Must(template.ParseGlob("form/*"))

func GetAll() *EmployeeList {
	query := "SELECT * FROM Employee"

	elist := new(EmployeeList)
	if err := queryRows(query, elist); err != nil {
		panic(err)
	}
	return elist
}

func GetById(id string) *Employee {
	query := "SELECT * FROM Employee WHERE id=?"
	emp := new(Employee)
	if err := queryRow(query, emp, id); err != nil {
		panic(err)
	}
	return emp
}

func Index(w http.ResponseWriter, r *http.Request) {
	res := GetAll().Items

	tmpl.ExecuteTemplate(w, "Index", res)
}

func Show(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	emp := GetById(id)
	tmpl.ExecuteTemplate(w, "Show", emp)
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	emp := GetById(id)

	tmpl.ExecuteTemplate(w, "Edit", emp)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")

		query := "INSERT INTO Employee(name, city) VALUES(?,?)"
		ExecuteQuery(query, name, city)

		log.Println("INSERT: Name: " + name + " | City: " + city)
	}
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("id")

		query := "UPDATE Employee SET name=?, city=? WHERE id=?"
		ExecuteQuery(query, name, city, id)

		log.Println("UPDATE: Name: " + name + " | City: " + city + " | Id: " + id)
	}
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	query := "DELETE FROM Employee WHERE id=?"
	ExecuteQuery(query, id)

	log.Println("DELETE")
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8000", nil)
}
