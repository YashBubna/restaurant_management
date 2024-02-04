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

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

type InvoiceViewFormat struct {
	Invoice_id       string
	Order_id         string
	Payment_status   *string
	Payment_method   string
	Payment_due_date time.Time
	Payment_due      interface{}
	Table_number     interface{}
	Order_details    interface{}
}

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		result, err := invoiceCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var allInvoices []bson.M
		err = result.All(ctx, &allInvoices)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		invoiceID := c.Param("invoice_id")

		var invoice models.Invoice

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceID}).Decode(&invoice)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var invoiceView InvoiceViewFormat

		allOrderItems, err := ItemsByOrder(invoice.Order_id)
		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date

		invoiceView.Payment_method = "null"
		if invoice.Payment_method != nil {
			invoiceView.Payment_method = *invoice.Payment_method
		}
		invoiceView.Payment_status = *&invoice.Payment_status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_details = allOrderItems[0]["order_details"]

		c.JSON(http.StatusOK, invoiceView)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice
		// var order models.Order

		invoiceID := c.Param("invoice_id")

		err := c.BindJSON(&invoice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.M{"invoice_id": invoiceID}

		var updateObject primitive.D

		// if invoice.Order_id != nil {
		// 	err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
		// 	defer cancel()
		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	updateObject = append(updateObject, bson.E{Key: "order_id", Value: invoice.Order_id})
		// }

		status := "PENDING"
		if invoice.Payment_status != nil {
			invoice.Payment_status = &status
		}

		if invoice.Payment_status != nil {
			updateObject = append(updateObject, bson.E{Key: "payment_status", Value: invoice.Payment_status})
		}

		if invoice.Payment_method != nil {
			updateObject = append(updateObject, bson.E{Key: "payment_method", Value: invoice.Payment_method})
		}

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{Key: "updated_at", Value: invoice.Updated_at})

		upsert := true

		option := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObject},
			},
			&option,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice
		var order models.Order

		err := c.BindJSON(&invoice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}

		err = orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		status := "PENDING"
		if invoice.Payment_status != nil {
			invoice.Payment_status = &status
		}

		invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()

		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		result, err := invoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
