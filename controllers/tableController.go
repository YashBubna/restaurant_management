package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/YashBubna/restaurant_management/database"
	"github.com/YashBubna/restaurant_management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		tableID := c.Param("table_id")
		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableID}).Decode(&table)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, table)
	}
}

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		result, err := tableCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var allTables []bson.M
		err = result.All(ctx, allTables)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, &allTables)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table

		err := c.BindJSON(&table)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(table)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, err := tableCollection.InsertOne(ctx, table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order item is not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table

		tableID := c.Param("table_id")
		filter := bson.M{"table_id": tableID}

		err := c.BindJSON(&table)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObject primitive.D

		if table.Table_no != nil {
			updateObject = append(updateObject, bson.E{Key: "table_no", Value: table.Table_no})
		}
		if table.Number_of_guests != nil {
			updateObject = append(updateObject, bson.E{Key: "number_of_guests", Value: table.Number_of_guests})
		}
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{Key: "updated_at", Value: table.Updated_at})

		upsert := true
		option := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObject},
			},
			&option,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
