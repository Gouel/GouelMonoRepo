package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    string             `bson:"firstName,omitempty"`
	LastName     string             `bson:"lastName,omitempty"`
	Email        string             `bson:"email,omitempty"`
	DOB          string             `bson:"dob,omitempty"` // Format YYYY-MM-DD
	Password     string             `bson:"password,omitempty"`
	Role         string             `bson:"role,omitempty"`
	Solde        float32            `bson:"solde"`
	Transactions []Transaction      `bson:"transactions,omitempty"`
}

type PurchaseProduct struct {
	ProductCode string `json:"productCode"`
	Amount      int    `json:"amount"`
}

type Transaction struct {
	Type         string            `bson:"type,omitempty"` // ex: "credit", "debit"
	Date         string            `bson:"date,omitempty"` // Format: "YYYY-MM-DDTHH:MM"
	EventId      string            `bson:"EventId,omitempty"`
	Cart         []PurchaseProduct `bson:"cart,omitempty"`
	Amount       float32           `bson:"amount,omitempty"`
	PayementType string            `bson:"payement_type,omitempty"` // "espece", "carte", "helloasso" ...
}
