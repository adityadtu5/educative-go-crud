package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	Routers()
}


func Routers() {
	InitDB()
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/users",
		GetUsers).Methods("GET")
	router.HandleFunc("/users",
		CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}",
		GetUser).Methods("GET")
	router.HandleFunc("/users/{id}",
		UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}",
		DeleteUser).Methods("DELETE")
	http.ListenAndServe(":3000",
		&CORSRouterDecorator{router})
}


// Task 3: Write your code here

// Task 4: write your code for create user here 

// Task 5: Write code for get user here

// Task 6: write code for update user here

// Task 7: Write code for delete user here

// Task 2: Write your code here

var db *sql.DB
var err error

// Task 1: Write code for DB initialization here

type CORSRouterDecorator struct {
	R *mux.Router
}

func (c *CORSRouterDecorator) ServeHTTP(rw http.ResponseWriter,
	req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Accept-Language,"+
				" Content-Type, YourOwnHeader")
	}

	if req.Method == "OPTIONS" {
		return
	}

	c.R.ServeHTTP(rw, req)
}

func InitDB(){
	db, err = sql.Open("mysql", "user:user@tcp(127.0.0.1:3306)/userdb")
	if err != nil {
		panic(err.Error())
	}
}