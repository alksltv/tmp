package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	IP   string
	Port string
	Res  string
}

var ip, port, chktcp string

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/", send).Methods("POST")
	r.HandleFunc("/result", result).Methods("GET")

	log.Print("Starting Server...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/home.html", nil)
}

func send(w http.ResponseWriter, r *http.Request) {
	// Step 1: Parse (and validate) form
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	ip = r.FormValue("ip")
	port = r.FormValue("port")

	// Step 2: Send message in an email
	chktcp = TCPClient(ip + ":" + port)

	// Step 3: Redirect to result page
	http.Redirect(w, r, "/result", http.StatusSeeOther)
}

func result(w http.ResponseWriter, r *http.Request) {
	msg := Message{
		IP:   ip,
		Port: port,
		Res:  chktcp,
	}
	fmt.Println(msg)
	render(w, "templates/result.html", msg)
}
func render(w http.ResponseWriter, filename string, data interface{}) {
	//fmt.Println(data)
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}

func TCPClient(id string) string {

	conn, err := net.DialTimeout("tcp", id, 7*time.Second)
	if err != nil && strings.Contains(err.Error(), "connection refused") {
		return fmt.Sprintf("TCP connection refused from: %s\n", id)

	} else if err != nil && strings.Contains(err.Error(), "timeout") {
		return fmt.Sprintf("TCP connection timed out: %s\n", id)

	} else if err != nil {
		return fmt.Sprintf("TCP connection failed: %s\n", id)

	} else {
		conn.Close()
		return fmt.Sprintf("TCP connected to: %s", id)
	}
}
