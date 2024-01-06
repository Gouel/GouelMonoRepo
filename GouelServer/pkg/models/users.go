package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"firstName,omitempty"`
	LastName  string             `bson:"lastName,omitempty"`
	Email     string             `bson:"email,omitempty"`
	DOB       string             `bson:"dob,omitempty"` // Format YYYY-MM-DD
	Password  string             `bson:"password,omitempty"`
	Role      string             `bson:"role,omitempty"`
	Solde     Solde              `bson:"solde,omitempty"`
}

type Solde struct {
	Amount       float64       `bson:"amount"`
	Transactions []Transaction `bson:"transactions,omitempty"`
}

type PurchaseItems []struct {
	ProductCode string `json:"product_code"`
	Amount      int    `json:"amount"`
}

type Transaction struct {
	Type         string        `bson:"type,omitempty"` // ex: "credit", "debit" , "refund"
	Date         string        `bson:"date,omitempty"` // Format: "YYYY-MM-DDTHH:MM"
	EventID      string        `bson:"event_id,omitempty"`
	Cart         PurchaseItems `bson:"cart,omitempty"`
	Amount       float64       `bson:"amount,omitempty"`
	PayementType string        `bson:"payement_type,omitempty"` // "espece", "carte", "helloasso" ...
}
