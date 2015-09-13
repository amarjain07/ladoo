package main

import (
	"fmt"
	"net/http"
	"html/template"
	"log"
	"os"
	"path"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/amarjain07/ladoo/gcm"
)

func main() {
	fmt.Println("Welcome")
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css")))) 
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img")))) 
	http.HandleFunc("/register/", registrationHandler)
	http.HandleFunc("/gcm/register/", gcmHandler)
	http.HandleFunc("/gcm/send/all/", gcmSendHandler)
	http.HandleFunc("/", templateHandler)
	http.ListenAndServe(":8080", nil)
}

func gcmSendHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method == "POST"){
		gcm.SendToAll(w, r)
	}else{
		http.NotFound(w, r)
	}
}

func gcmHandler(w http.ResponseWriter, r *http.Request) {
	if(r.Method == "POST"){
		gcm.Register(w, r)
	}else{
		http.NotFound(w, r)
	}
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "/test")
	if err != nil {
    	    
    }
    defer db.Close()
    query := "INSERT INTO registration (name, mobileNo, couponNo, imei, dealerName) VALUES(? ,? , ? , ? , ?)"
    stmtIns, err := db.Prepare(query)
    if err != nil {
    	panic(err.Error())
    }
    defer stmtIns.Close()

    _, err = stmtIns.Exec(r.FormValue("name"), r.FormValue("mobile"), r.FormValue("coupon"), r.FormValue("imei"),r.FormValue("dealer"))
    if err != nil {
    	panic(err.Error())
    }

    fp := path.Join("templates", "register.html")
    tmpl,err := template.ParseFiles(fp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	if urlPath == "/" {
		urlPath = "index.html"
	}
	fp := path.Join("templates", urlPath)
	info, err := os.Stat(fp)
	if  err != nil{
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl,err := template.ParseFiles(fp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}