package model

type Transfer struct {
	ID                   int
	SenderID, ReceiverID int
	Amount               int
}
