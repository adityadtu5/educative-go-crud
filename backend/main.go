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

type User struct {
	ID          string `json:"id"`
	FirstName   string `json:"firstName"`
	MiddleName  string `json:"middleName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	CivilStatus string `json:"civilStatus"`
	BirthDay    string `json:"birthday"`
	Contact     string `json:"contact"`
	Age         string `json:"address"`
	Address     string `json:"age"`
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

// Task 7: Write code for delete user here


var db *sql.DB
var err error


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

func InitDB() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/userdb")
	if err != nil {
		panic(err.Error())
	}
}

func GetUsers(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var users []User
	result, err := db.Query("SELECT `id`, `first_name`, `middle_name`, `last_name`, `email`, `gender`, `civil_status`, `birthday`, `contact`, `address`, floor(datediff(now(),birthday)/365) AS age FROM users")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer result.Close()
	for result.Next() {
		var user User
		err := result.Scan(
			&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email,
			&user.Gender, &user.CivilStatus, &user.BirthDay, &user.Contact, &user.Address, &user.Age)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(users)
}

func CreateUser(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "applicatio/json")
	
	stmt, err := db.Prepare("INSERT into users (first_name, middle_name, last_name, email, gender, civil_status, birthday, contact, address) values(?,?,?,?,?,?,?,?,?)")
	if err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var keyVal map[string]string
	err = json.Unmarshal(body, &keyVal)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	firstName := keyVal["firstName"]
	middleName := keyVal["middleName"]
	lastName := keyVal["lastName"]
	email := keyVal["email"]
	gender := keyVal["gender"]
	civilStatus := keyVal["civilStatus"]
	birthDay := keyVal["birthday"]
	contact := keyVal["contact"]
	address := keyVal["address"]

	_, err = stmt.Exec(firstName, middleName, lastName, email, gender, civilStatus, birthDay, contact, address)

	if err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "New user was created")
}

func GetUser(rw http.ResponseWriter, req *http.Request){
	rw.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	result, err := db.Query("SELECT id, first_name, middle_name, last_name, email, gender, civil_status, birthday, contact, address, floor(datediff(now(), birthday)/365) as age FROM users WHERE id = ?", params["id"])
	if err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer result.Close()
	var user User
	userFound := false
	for result.Next(){
		err := result.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.Gender, &user.CivilStatus, &user.BirthDay, &user.Contact, &user.Age, &user.Address )
		if err != nil{
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		userFound = true
	}
	if err := result.Err(); err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if !userFound{
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "User not found with ID: %s", params["id"])
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(user)
}

func UpdateUser(rw http.ResponseWriter, req *http.Request){
	rw.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	stmt, err := db.Prepare("UPDATE users SET first_name = ?, middle_name = ?, last_name = ?, email = ?, gender = ? civil_status = ?, birthday = ?, contact = ?, address = ? WHERE id=?")
	if err != nil{
		panic(err.Error())
	}
	defer stmt.Close()
	var userUpdate User
	if err := json.NewDecoder(req.Body).Decode(&userUpdate); err != nil{
		panic(err.Error())
	}
	result, err := stmt.Exec(userUpdate.FirstName, userUpdate.MiddleName, userUpdate.LastName, userUpdate.Email, userUpdate.Gender, userUpdate.CivilStatus, userUpdate.BirthDay, userUpdate.Contact, userUpdate.Address, params["id"])
	if err != nil {
		panic(err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		panic(err.Error())
	}
	if rowsAffected == 0 {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "No user found with id = %s", params["id"])
		return
	}
	fmt.Fprintf(rw, "User with ID = %s was updated", params["id"])
}

func DeleteUser(rw http.ResponseWriter, req *http.Request){
	rw.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	stmt, err := db.Prepare("DELETE from users where id = ?")
	if err != nil{
		panic(err.Error())
	}
	defer stmt.Close()
	result, err := stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		panic(err.Error())
	}
	if rowsAffected == 0 {
		http.Error(rw, fmt.Sprintf("No user found with id = %s", params["id"]), http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "User with ID = %s was deleted", params["id"])
}