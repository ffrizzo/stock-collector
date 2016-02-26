package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ffrizzo/stock-collector/models"
	"github.com/gin-gonic/gin"
)

//CollectorController represents the controller for Colletor resource
type CollectorController struct {
	db *sql.DB
}

//NewCollectorController instantiate a new controller for collectors
func NewCollectorController(db *sql.DB) *CollectorController {
	return &CollectorController{db}
}

//Collector retrieve the data, validate and save
func (controller CollectorController) Collector(c *gin.Context) {
	var collector models.Stock
	c.BindJSON(&collector)

	if collector.Account == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Account cannot be null."})
		return
	}

	if collector.User == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "User cannot be null."})
		return
	}

	if collector.Ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Ticker cannot be null."})
		return
	}

	if collector.Time == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Time cannot be null."})
		return
	}

	tx := c.MustGet("tx").(*sql.Tx)
	stock, err := SaveStock(tx, collector)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": true, "message": fmt.Sprintf("Error to save data for stock. Error: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "stock": stock})
}

//SaveStock save data for stock
func SaveStock(tx *sql.Tx, collector models.Stock) (*models.Stock, error) {
	query := `insert into stock (sell, rate, buy, ticker, account, username, time)
        values ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	var stockID int64
	err := tx.QueryRow(query, collector.Sell, collector.Rate, collector.Buy, collector.Ticker,
		collector.Account, collector.User, collector.Time).Scan(&stockID)
	if err != nil {
		return nil, err
	}

	stock := &models.Stock{stockID, collector.Sell, collector.Rate, collector.Buy, collector.Ticker,
		collector.Account, collector.User, collector.Time}

	return stock, nil
}
