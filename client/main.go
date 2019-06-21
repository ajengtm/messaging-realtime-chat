package main

/*
@author ajeng tya
@date 10 Juni 2019

Microservice Massaging as Web Socket Client
*/

import (
	"bitbucket.org/kudoindonesia/messaging-realtime-chat/client/pkg/responses"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// 0 pendinf 1 delivered, 2 read

type Message struct {
	Id 				int	`json:"sender_id"`
	SenderId 		int	`json:"sender_id"`
	RecipientId    int `json:"recepient_id"`
	Message  	string `json:"message"`
	Status 		int	`json:"sender_id"`
	CreatedAt   string
	UpdatedAt   string
}

// CreateMessageParam for parameter CreateMessage
type CreateMessage struct {
	SenderId 		int	`json:"sender_id"`
	RecipientId    int `json:"recepient_id"`
	Message  string `json:"message"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "messaging"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	//db, err := sql.Open("mysql", "root:''@tcp(127.0.0.1:3306)/messaging")
	if err != nil {
		panic(err.Error())
	}
	return db
}

// API for collect message that has been sent out
func getMessage(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM messages ORDER BY id DESC ")
	if err != nil {
		panic(err.Error())
	}
	var msg Message
	for selDB.Next() {
		err := selDB.Scan(
			&msg.Id,
			&msg.SenderId,
			&msg.RecipientId,
			&msg.Message,
			&msg.Status,
			&msg.UpdatedAt,
			&msg.CreatedAt,
		)
		if err != nil {
			log.Println(err)
		}
		log.Println(msg)
	}
	defer db.Close()

	responses.Write(w, responses.APICreated)
	return
}

func saveMessage(senderId int, recipientId int, msg string)(err error){
	var message Message

	message.SenderId = senderId
	message.RecipientId = recipientId
	message.Message = msg

	//message status delivered
	message.Status = 2

	current_time := time.Now().Local()
	now := current_time.Format("2006-01-02 15:04:05")
	message.UpdatedAt = now
	message.CreatedAt = now

	log.Println(message)
	db := dbConn()

	stmtIns, err := db.Prepare("INSERT INTO messages (sender_id, recepient_id, message, status, updated_at, created_at) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	log.Println(stmtIns)

	res, err := stmtIns.Exec(
		message.SenderId,
		message.RecipientId,
		message.Message,
		message.Status,
		message.UpdatedAt,
		message.CreatedAt,
	)
	defer db.Close()
	return err

	lastID, err := res.LastInsertId()
	log.Println(lastID)
	return err
}

var addr = flag.String("addr", "localhost:8080", "http service address")
// API for sending a message
func sendMessage(w http.ResponseWriter, r *http.Request) {

	senderId, _ := strconv.Atoi( string(r.FormValue("sender_id")))
	recipientId, _ := strconv.Atoi( string(r.FormValue("recepient_id")))
	var message = string(r.FormValue("message"))

	msg := &CreateMessage{senderId, recipientId,message}
	msgWs, err := json.Marshal(msg)
	if err != nil {
		panic (err)
	}

	fmt.Println(string(msgWs))


	//start process websocket
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

		// send message to websocket server
	err = c.WriteMessage(websocket.TextMessage, []byte(msgWs))
	if err != nil {
		log.Println("write:", err)
		return
	}
	// end proses websocket

	err = saveMessage(senderId, recipientId, message)
	if err != nil {
		log.Println("write:", err)
		return
	}

	responses.Write(w, responses.APICreated)
	return

}

func main() {

	http.HandleFunc("/message", sendMessage)
	http.HandleFunc("/message/history", getMessage)

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}


}
