package main

import (
	"github.com/dgrijalva/jwt-go"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/url"
	"strings"
	"fmt"
)

type JWTResponse struct {
	ConsumerID string `json:"consumer_id"`
	ID string `json:"id"`
	Secret string `json:"secret"`
	Key string `json:"key"`
	CreatedAt int64 `json:"created_at"`
	Algorithm string `json:"algorithm"`
}

type SignUp struct {
	Login string `json:"login"`
	Pass string `json:"pass"`
}

type SignUpConsumer struct {
	Username string `json:"username"`
	CreatedAt int64 `json:"created_at"`
	ID string `json:"id"`
}

type SignUpBasicAuth struct {
	Password string `json:"password"`
	ConsumerID string `json:"consumer_id"`
	ID string `json:"id"`
	Username string `json:"username"`
	CreatedAt int64 `json:"created_at"`
}

func craftJWT(jwtObj JWTResponse) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": jwtObj.Key,
	})

	tokenString, err := token.SignedString([]byte(jwtObj.Secret))

	if err != nil {
		panic(err)
	}

	return tokenString
}

func getJWT(login string) JWTResponse {
	endpoint := "http://oscars_kong:8001/consumers/" + login + "/jwt"

	contentReader := bytes.NewReader([]byte("{}"))
	req, _ := http.NewRequest("POST", endpoint, contentReader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	htmlData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var jtwObj JWTResponse
	err2 := json.Unmarshal(htmlData, &jtwObj);

	if err2 != nil {
		panic(err2)
	}

	return jtwObj
}

func getToken(login string) string {

	if (len(login) <= 0) {
		return ""
	}

	jwtObj := getJWT(login)
	JWTString := craftJWT(jwtObj)

	return JWTString
}

func createBasicAuth(login, pass string) bool {

	endpoint := "http://oscars_kong:8001/consumers/"+ login +"/basic-auth"

	form := url.Values{}
	form.Add("username", login)
	form.Add("password", pass)
	hc := http.Client{}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := hc.Do(req)

	if err != nil {
		panic(err)
	}

	htmlData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var signUpKong SignUpBasicAuth
	err2 := json.Unmarshal(htmlData, &signUpKong);

	if err2 != nil {
		return false
	}

	if len(signUpKong.ID) <= 0 {
		return false
	}

	return true
}

func createUser(login string) bool {

	endpoint := "http://oscars_kong:8001/consumers"

	form := url.Values{}
	form.Add("username", login)
	hc := http.Client{}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := hc.Do(req)

	if err != nil {
		panic(err)
	}

	htmlData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var signupkong SignUpConsumer
	err2 := json.Unmarshal(htmlData, &signupkong);

	if err2 != nil {
		return false
	}

	if len(signupkong.ID) <= 0 {
		return false
	}

	return true
}

func createKong(login, pass string) bool {

	success := createUser(login)

	if success == false {
		return false
	}

	success = createBasicAuth(login, pass)

	if success == false {
		return false
	}

	return true
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	login := r.Header.Get("X-Credential-Username")
	token := getToken(login)

	mapD := map[string]string{"token": token}
	mapB, err := json.Marshal(mapD)
	if err != nil {
		panic(err)
	}

	w.Write(mapB)
}

func serveCreate(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var signup SignUp
	err := decoder.Decode(&signup)

	fmt.Println(signup)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	success := createKong(signup.Login, signup.Pass)

	mapD := map[string]bool{"success": success}
	mapB, err := json.Marshal(mapD)
	if err != nil {
		panic(err)
	}

	w.Write(mapB)
}

func main () {
	r := mux.NewRouter()
	r.HandleFunc("/create", serveCreate)
	r.HandleFunc("/login", serveLogin)
	http.Handle("/", r)
	http.ListenAndServe(":80", nil)
}
