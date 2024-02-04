package controllers

import (
	"context"
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
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		var allMenus []bson.M
		if err = result.All(ctx, allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		menuID := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuID}).Decode(&menu)
		cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		err := c.BindJSON(&menu)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validattionErr := validate.Struct(menu)
		if validattionErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, err := menuCollection.InsertOne(ctx, menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu item is not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
		defer cancel()
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		err := c.BindJSON(&menu)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		menuID := c.Param("menu_id")
		filter := bson.M{"menu_id": menuID}

		var updateObject primitive.D

		if menu.Start_date != nil && menu.Start_date != nil {
			if !inTimeSpan(*menu.Start_date, *menu.End_date, time.Now()) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "kindly retype the time"})
				return
			}

			updateObject = append(updateObject, bson.E{Key: "start_date", Value: menu.Start_date})
			updateObject = append(updateObject, bson.E{Key: "end_date", Value: menu.End_date})

			if menu.Name != "" {
				updateObject = append(updateObject, bson.E{Key: "name", Value: menu.Name})
			}
			if menu.Category != "" {
				updateObject = append(updateObject, bson.E{Key: "category", Value: menu.Category})
			}

			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObject = append(updateObject, bson.E{Key: "updated_at", Value: menu.Updated_at})

			upsert := true

			option := options.UpdateOptions{
				Upsert: &upsert,
			}
			result, err := menuCollection.UpdateOne(
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
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}
