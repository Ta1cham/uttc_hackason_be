package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"uttc_hackason_be/controller"
	"uttc_hackason_be/dao"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPwd := os.Getenv("MYSQL_PWD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlPwd, mysqlHost, mysqlDatabase)
	var err error
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}
}

//func init() {
//	mysqlUser := "user"
//	mysqlUserpsw := "password"
//	mysqlDatabase := "mydatabase"
//
//	_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(localhost:3306)/%s", mysqlUser, mysqlUserpsw, mysqlDatabase))
//	if err != nil {
//		log.Fatalf("fail: sql.Open, %v\n", err)
//	}
//	if err := _db.Ping(); err != nil {
//		log.Fatalf("fail: _db.Ping, %v\n", err)
//	}
//	db = _db
//}

func handler(w http.ResponseWriter, r *http.Request) {
	userdao := dao.CreateDao(db)
	switch r.Method {
	case http.MethodGet:
		controller.GetUserController(w, r, userdao)
	case http.MethodPost:
		controller.RegisterUserController(w, r, userdao)
	default:
		log.Printf("fail: HTTP Method is %s\n", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	http.HandleFunc("/user", handler)
	closeDBWithSysCall()

	log.Println("Listening...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func closeDBWithSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received syscall, %v", s)

		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("success: db.Close()")
		os.Exit(0)
	}()
}
