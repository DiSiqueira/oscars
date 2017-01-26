//app.go

package main

import "net/http"
import "html/template"

type Variables struct{
	Title string
	Heading string
}

var templates = template.Must(template.ParseFiles("oscars.html"))

func serve(res http.ResponseWriter, req *http.Request){

	myVars := Variables{"My Website Title", "My Website Heading"}

	templates.ExecuteTemplate(res, "oscars.html", myVars)

}

func main(){

	http.HandleFunc("/", serve)
	http.ListenAndServe(":3000", nil)
}