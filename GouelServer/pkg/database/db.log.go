package database

import (
	"time"

	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func LogAction(c *gin.Context, data interface{}) {

	userID, _ := c.Get("userId")
	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return
	}

	log := models.AccessLog{
		UserID:    userOID,
		Timestamp: time.Now(),
		Route:     c.Request.URL.Path,
		Method:    c.Request.Method,
		Data:      data,
	}

	if c.Request.Method != "GET" {
		collection := Database.Collection("log")
		collection.InsertOne(c, log)
	}

}
