package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)
// Define our message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

var addr = flag.String("addr", "localhost:8080", "http service address")

func sendMessage(w http.ResponseWriter, r *http.Request) {

	var email = string(r.FormValue("email"))
	var username = string(r.FormValue("username"))
	var message = string(r.FormValue("message"))

	//ini proses websocket nya
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	// end proses websocket

	a := &Message{email, username, message}

	out, err := json.Marshal(a)
	if err != nil {
		panic (err)
	}

	fmt.Println(string(out))

	// kirim ke server websocket
	err = c.WriteMessage(websocket.TextMessage, []byte(out))
	if err != nil {
		log.Println("write:", err)
		return
	}

	fmt.Fprintf(w, message)
	return

}

func main() {
	http.HandleFunc("/message", sendMessage)

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}


}
