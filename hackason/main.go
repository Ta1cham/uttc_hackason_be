package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"uttc_hackason_be/controller"
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

func main() {
	r := mux.NewRouter()

	userController := controller.NewUserController(db)
	tweetController := controller.NewTweetController(db)
	likesController := controller.NewLikesController(db)
	noteController := controller.NewNoteController(db)
	userController.RegiterRoutes(r)
	tweetController.RegisterRoute(r)
	likesController.RegisterRoute(r)
	noteController.RegiterRoutes(r)

	closeDBWithSysCall()

	c := cors.Default()
	handler := c.Handler(r)

	log.Println("Listening...")
	if err := http.ListenAndServe(":8000", handler); err != nil {
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
