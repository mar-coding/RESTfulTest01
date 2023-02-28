package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var counter int
var mutex = &sync.Mutex{}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	fmt.Println("Endpoint Hit: homePage")
}

func aminPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi")
	fmt.Println("Endpoint Hit: aminPage")
}

func incrementCounterPage(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, strconv.Itoa(counter))
	mutex.Unlock()
	fmt.Println("Endpoint Hit: incPage")
}

func handleRequests() {

	http.HandleFunc("/amin", aminPage)
	http.HandleFunc("/inc", incrementCounterPage)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	// http.HandleFunc("/", indexPage)

	//log.Fatal(http.ListenAndServe(":12345", nil))
	log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))
}

func main() {
	handleRequests()
}
