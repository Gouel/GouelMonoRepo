package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gouel/gouel_serveur/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(user models.User) (string, error) {
	collection := Database.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.Solde = models.Solde{Amount: 0.0, Transactions: []models.Transaction{}}

	if user.Password != "" {
		user.Password = HashPassword(user.Password)
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func VerifyUser(email, password string) (models.User, error) {
	collection := Database.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("couple email/password invalide") // Utilisateur non trouvé ou erreur
	}

	if HashPassword(password) != user.Password {
		return models.User{}, fmt.Errorf("couple email/password invalide") // Le mot de passe ne correspond pas
	}

	return user, nil // Le couple email/mot de passe est valide
}

func UpdateUser(userID string, updateData bson.M) error {
	collection := Database.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Vérifier si le mot de passe est présent dans les données de mise à jour
	if password, ok := updateData["password"].(string); ok && password != "" {

		updateData["password"] = HashPassword(password)
	}

	objID, _ := primitive.ObjectIDFromHex(userID)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	return err
}

func UpdateUserBalance(userID string, amount float64) error {
	collection := Database.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(userID)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"solde.amount": amount}})
	return err
}

func AddUserTransaction(userID string, transaction models.Transaction) error {
	collection := Database.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(userID)

	// Récupérer l'utilisateur et son solde actuel
	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return err
	}

	newBalance := user.Solde.Amount

	switch transaction.Type {
	case "credit":
		newBalance += transaction.Amount
	case "debit", "refund":
		if newBalance < transaction.Amount {
			return fmt.Errorf("solde insuffisant")
		}
		newBalance -= transaction.Amount
	default:
		return fmt.Errorf("type de transaction non reconnu")
	}

	// Mettre à jour le solde de l'utilisateur
	update := bson.M{
		"$set":  bson.M{"solde.amount": newBalance},
		"$push": bson.M{"solde.transactions": transaction},
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func GetUserByID(userID string) (bson.M, error) {
	collection := Database.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(userID)

	var user bson.M
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func FindUsersByEmailStartsWith(emailStart string) ([]bson.M, error) {
	collection := Database.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"email": bson.M{"$regex": "^" + emailStart, "$options": "i"}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []bson.M
	for cursor.Next(ctx) {
		var user bson.M
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		// Convertir l'ID en hexadécimal
		if oid, ok := user["_id"].(primitive.ObjectID); ok {
			user["_id"] = oid.Hex()
		}
		users = append(users, user)
	}

	return users, nil
}

func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
