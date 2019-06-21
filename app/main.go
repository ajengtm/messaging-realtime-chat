package main

/*
@author ajeng tya
@date 10 Juni 2019

Microservice Massaging as Web Socket Client
*/

import (
	"database/sql"
	"encoding/json"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)


// APIResponse defines attributes for api Response
type APIResponse struct {
	HTTPCode   int         `json:"-"`
	Code       int         `json:"code"`
	Message    interface{} `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

// Messages defines attributes for record message in database
type Messages struct {
	Id 				int	`json:"id"`
	SenderId 		int	`json:"sender_id"`
	RecipientId    	int `json:"recepient_id"`
	Message  		string 	`json:"message"`
	Status 			int		`json:"status"`
	CreatedAt   	string	`json:"created_at"`
	UpdatedAt   	string	`json:"updated_at"`
}

// MessagesResponse defines attributes for api response get Messages
type MessagesResponse struct {
	Id 				int	`json:"id"`
	SenderId 		int	`json:"sender_id"`
	RecipientId    	int `json:"recepient_id"`
	Message  		string 	`json:"message"`
	Status 			int		`json:"status"`
	StatusName 		string	`json:"status_name"`
	CreatedAt   	string	`json:"created_at"`
	UpdatedAt   	string	`json:"updated_at"`
}

// CreateMessageParam for parameter CreateMessage to websocket
type CreateMessage struct {
	SenderId 		int	`json:"sender_id"`
	RecipientId    	int `json:"recepient_id"`
	Message  		string `json:"message"`
}

var (
	APICreated = APIResponse{
		HTTPCode: http.StatusCreated,
		Code:     1000,
		Message:  "Success",
	}
	addr = flag.String("addr", "localhost:8080", "http service address")
)

func main() {

	http.HandleFunc("/message", sendMessage)
	http.HandleFunc("/message/history", getMessage)

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/* Connection to Database */
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "messaging"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

/* Get status message */
func getStatusName(status int) (name string){
	switch status {
		case 0:
			name = "Pending"
		case 1:
			name = "Delivered"
		case 2:
			name = "Read"
		default:
			name = "Undefined"
		}
	return name
}

/* API for collect message that has been sent out by sender_id */
func getMessage(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	senderId, _ := strconv.Atoi( string(r.FormValue("sender_id")))

	query := `SELECT id, sender_id, recepient_id, message, status, updated_at, created_at
			  FROM messages 
			  WHERE sender_id = ? 
			  ORDER BY id DESC `
	selDB, err := db.Query(query, senderId)
	if err != nil {
		panic(err.Error())
	}
	var a []MessagesResponse
	for selDB.Next() {
		message := MessagesResponse{}
		err := selDB.Scan(
			&message.Id,
			&message.SenderId,
			&message.RecipientId,
			&message.Message,
			&message.Status,
			&message.UpdatedAt,
			&message.CreatedAt,
		)

		// skip when scan error
		if err != nil {
			log.Println(err)
			continue
		}

		message.StatusName= getStatusName(message.Status)
		a = append(a, message)

	}
	defer db.Close()

	Write(w, APICreated.WithData(a))
	return
}

/* Function for save massage to database , part of sending message process*/
func saveMessage(senderId int, recipientId int, msg string)(err error){
	var message Messages

	message.SenderId = senderId
	message.RecipientId = recipientId
	message.Message = msg

	//message status delivered
	message.Status = 1

	current_time := time.Now().Local()
	now := current_time.Format("2006-01-02 15:04:05")
	message.UpdatedAt = now
	message.CreatedAt = now

	db := dbConn()

	stmtIns, err := db.Prepare("INSERT INTO messages (sender_id, recepient_id, message, status, updated_at, created_at) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmtIns.Exec(
		message.SenderId,
		message.RecipientId,
		message.Message,
		message.Status,
		message.UpdatedAt,
		message.CreatedAt,
	)

	defer db.Close()
	return err
}

/* Function for sending message to websocket */
func sendMessage(w http.ResponseWriter, r *http.Request) {

	senderId, _ := strconv.Atoi( string(r.FormValue("sender_id")))
	recipientId, _ := strconv.Atoi( string(r.FormValue("recepient_id")))
	message := string(r.FormValue("message"))

	msg := &CreateMessage{senderId, recipientId,message}
	msgWs, err := json.Marshal(msg)
	if err != nil {
		panic (err)
	}

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

	// save massage to database
	err = saveMessage(senderId, recipientId, message)
	if err != nil {
		log.Println("writeBD:", err)
		return
	}

	Write(w, APICreated.WithMessage("Delivered"))
	return

}


/* Format Response API */
func Write(w http.ResponseWriter, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.HTTPCode)
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(js); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *APIResponse) WithData(data interface{}) APIResponse {
	new := new(APIResponse)
	new.HTTPCode = a.HTTPCode
	new.Code = a.Code
	new.Message = a.Message
	new.Data = data

	return *new
}

func (a *APIResponse) WithMessage(message interface{}) APIResponse {
	new := new(APIResponse)
	new.HTTPCode = a.HTTPCode
	new.Code = a.Code
	new.Message = message
	new.Data = a.Data

	return *new
}
/* ****************** */