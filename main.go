package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="userName" name="userName">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Cookie)
	name := getUserName(request)
	if name != "" {
		fmt.Fprintf(response, internalPage, name)
	} else {
		http.Redirect(response, request, "/", 302)
	}

}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("userName")
	pass := request.FormValue("password")
	fmt.Println(name)
	redirectTarget := "/"
	if name != "" && pass != "" {
		// .. check credentials ..
		setSession(name, response)
		redirectTarget = "/internal"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

var router = mux.NewRouter()

func main() {

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/internal", internalPageHandler)

	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/logout", logoutHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
