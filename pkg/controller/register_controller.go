package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twilio-labs/sample-twilio-go/pkg/configuration"
	"github.com/twilio-labs/sample-twilio-go/pkg/db"
)

type RegisterController struct {
	ctx context.Context
	db  *db.DB
}

func NewRegisterController(ctx context.Context, db *db.DB) *RegisterController {
	return &RegisterController{ctx, db}
}

func (ctr *RegisterController) GET(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func (ctr *RegisterController) POST(c *gin.Context) {
	var customer configuration.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	customer.ID = configuration.GenerateUserID()
	customer.CreatedAt = fmt.Sprint(time.Now().UnixMilli())

	err := ctr.db.CreateCustomer(ctr.ctx, &customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "internalError",
			"message": "An internal error has occurred.",
		})
	}

	json, _ := json.Marshal(configuration.CustomerToMap(&customer))
	fmt.Println("Customer: " + string(json))
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Success",
		"customer": string(json),
	})
}
