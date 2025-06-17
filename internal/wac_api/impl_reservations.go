package wac_api

import (
	"net/http"
	"github.com/SimonValicek/wac-project-api/internal/database"
	"github.com/SimonValicek/wac-project-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// ✅ This is your struct implementing the DefaultAPI interface
type implReservationAPI struct{}

// ✅ Constructor that returns the interface implementation
func NewReservationApi() DefaultAPI {
	return &implReservationAPI{}
}

// ✅ These methods implement the DefaultAPI interface

func (o *implReservationAPI) ReservationsGet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.Collection.Find(ctx, bson.M{})
   if err != nil {
      // expose the underlying error for debugging
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
	}
	defer cursor.Close(ctx)

	var reservations []models.Reservation
	if err := cursor.All(ctx, &reservations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode results"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}



func (o *implReservationAPI) ReservationsPost(c *gin.Context) {
	var input models.Reservation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := database.Collection.InsertOne(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert reservation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"inserted_id": res.InsertedID})
}


func (o *implReservationAPI) ReservationsIdPut(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	var updated models.Reservation
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid object id"})
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
    "$set": bson.M{
        "licensePlate": updated.LicensePlate,
        "category":     updated.Category,
        "datetime":     updated.Datetime,
        "spotNumber":   updated.SpotNumber,
    },
}


	result, err := database.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

func (o *implReservationAPI) ReservationsIdDelete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid object id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := database.Collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

