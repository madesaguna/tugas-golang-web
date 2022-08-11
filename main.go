package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id       int
	TaskName string
	Receiver string
	DateLine string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "tanyagl3nn"
	dbName := "golangweb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("views/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT id, taskname, receiver, dateline FROM task ORDER BY dateline DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var taskName, receiver, dateLine string
		err = selDB.Scan(&id, &taskName, &receiver, &dateLine)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.TaskName = taskName
		emp.Receiver = receiver
		emp.DateLine = dateLine
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
	//log.Print(time.F)
}

func View(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT id, taskname, receiver, dateline FROM task WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var taskName, receiver, dateLine string
		err = selDB.Scan(&id, &taskName, &receiver, &dateLine)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.TaskName = taskName
		emp.Receiver = receiver
		emp.DateLine = dateLine
	}
	tmpl.ExecuteTemplate(w, "View", emp)
	defer db.Close()
}

func Create(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "Create", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT id, taskname, receiver, dateline FROM task WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var taskName, receiver, dateLine string
		err = selDB.Scan(&id, &taskName, &receiver, &dateLine)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.TaskName = taskName
		emp.Receiver = receiver
		emp.DateLine = dateLine
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Save(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Print('1')
	if r.Method == "POST" {
		taskName := r.FormValue("taskname")
		receiver := r.FormValue("receiver")
		dateLine := r.FormValue("dateline")
		insForm, err := db.Prepare("INSERT INTO task(taskname, receiver, dateline) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(taskName, receiver, dateLine)
		log.Println("INSERT: Name: " + taskName + " | City: " + receiver)
	}
	defer db.Close()
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		taskName := r.FormValue("taskname")
		receiver := r.FormValue("receiver")
		dateLine := r.FormValue("dateline")
		id := r.FormValue("id")
		insForm, err := db.Prepare("UPDATE task SET taskname=?, receiver=?, dateline=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(taskName, receiver, dateLine, id)
		log.Println("UPDATE: Name: " + taskName + " | City: " + receiver)
	}
	defer db.Close()
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM task WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/view", View)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/save", Save)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
