/*
	PROGRAM DOESN'T FULLY WORK FOR NOW
	PROGRAM DOESN'T FULLY WORK FOR NOW
	PROGRAM DOESN'T FULLY WORK FOR NOW
	PROGRAM DOESN'T FULLY WORK FOR NOW
	PROGRAM DOESN'T FULLY WORK FOR NOW
*/

package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

// LogIn dupa
type LogIn struct {
	Username string
	Password string
}

// User dupa
type User struct {
	id        int
	username  string
	password  string
	createdAt time.Time
}

var createTableQuery = `
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME,
    PRIMARY KEY (id)
);
`

func dbInitConn() *sql.DB {
	db, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/gobank?parseTime=true")
	if err != nil {
		log.Error("Unsuccessful connection open with database, ", err)
	}
	return db
}

var db = dbInitConn()

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// Account dupa
func Account(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "dupa")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	println(w, "The cake is a lie!")
}

// Login dupa
func Login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./site/account.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	session, _ := store.Get(r, "dupa")
	/*details := LogIn{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}*/

	rows, err := db.Query(`SELECT username, password FROM users`)
	if err != nil {
		log.Warn("cos tam chuj")
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.username, &u.password)
		if err != nil {
			log.Warn("cos tam dupa")
		}
		users = append(users, u)
	}
	err = rows.Err()
	println(users)
	session.Values["authenticated"] = true
	session.Save(r, w)

	tmpl.Execute(w, struct{ Success bool }{true})
}

// Logout dupa
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "dupa")

	session.Values["authenticated"] = false
	session.Save(r, w)
}

func main() {
	// ale tu burdel
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	log.SetOutput(os.Stdout)
	// TODO: add hashing password with golang.org/x/crypto/bcrypt

	defer db.Close()
	err := db.Ping()
	if err != nil {
		log.Error("Database ping failure, ", err)
	}
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Error("dupaa fail, ", err)
	}
	fs := http.FileServer(http.Dir("site/"))
	http.Handle("/", http.StripPrefix("/", fs))
	http.HandleFunc("/secret", Account)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)

	http.ListenAndServe(":8080", nil)
}
