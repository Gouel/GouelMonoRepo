package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    string             `bson:"FirstName,omitempty"`
	LastName     string             `bson:"LastName,omitempty"`
	Email        string             `bson:"Email,omitempty"`
	DOB          string             `bson:"DOB,omitempty"` // Format YYYY-MM-DD
	Password     string             `bson:"Password,omitempty"`
	Role         string             `bson:"Role,omitempty"`
	Solde        float32            `bson:"Solde"`
	Transactions []Transaction      `bson:"Transactions,omitempty"`
}

type PurchaseProduct struct {
	ProductCode string `json:"ProductCode"`
	Amount      int    `json:"Amount"`
}

type Transaction struct {
	Type        string             `bson:"Type,omitempty"` // ex: "credit", "debit"
	Date        string             `bson:"Date,omitempty"` // Format: "YYYY-MM-DDTHH:MM"
	EventId     primitive.ObjectID `bson:"EventId,omitempty"`
	Cart        []PurchaseProduct  `bson:"Cart,omitempty"`
	Amount      float32            `bson:"Amount,omitempty"`
	PaymentType string             `bson:"PaymentType,omitempty"` // "espece", "carte", "helloasso" ...
}
