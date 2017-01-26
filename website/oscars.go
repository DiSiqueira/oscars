//app.go

package main

import "net/http"
import (
	"html/template"
	"github.com/gorilla/mux"
)

type Variables struct{
	Title string
	Heading string
}

var templates = template.Must(template.ParseFiles("index.html"))

func serve(res http.ResponseWriter, req *http.Request){

	myVars := Variables{"My Website Title", "My Website Heading"}

	templates.ExecuteTemplate(res, "index.html", myVars)

}

func main () {
	r := mux.NewRouter()
	r.HandleFunc("/", serve)

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/", r)
	http.ListenAndServe(":80", nil)
}
