package models

// Define Model of message
type Message struct {
	Id 				int		`json:"sender_id"`
	SenderId 		int		`json:"sender_id"`
	RecipientId    	int 	`json:"recepient_id"`
	Message  		string 	`json:"message"`
	Status 			int		`json:"sender_id"`
	CreatedAt   	string	`json:"created_at"`
	UpdatedAt   	string	`json:"updated_at"`
}
