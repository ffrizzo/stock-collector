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
		panic("Account cannot be null.")
	}

	if collector.User == "" {
		panic("User cannot be null.")
	}

	if collector.Ticker == "" {
		panic("Ticker cannot be null.")
	}

	if collector.Time == nil {
		panic("Time cannot be null.")
	}

	tx := c.MustGet("tx").(*sql.Tx)
	stock, err := SaveStock(tx, collector)
	if err != nil {
		panic(fmt.Sprintf("Error to save data for stock. Error: %s", err))
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "stock": stock})
}

//SaveStock save or update data for Stock
func SaveStock(tx *sql.Tx, collector models.Stock) (models.Stock, error) {
	stock := models.Stock{}

	query := `insert into stock (sell, rate, buy, ticker, account, username, time)
        values ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	var stockID int64
	err := tx.QueryRow(query, collector.Sell, collector.Rate, collector.Buy, collector.Ticker,
		collector.Account, collector.User, collector.Time).Scan(&stockID)
	if err != nil {
		return models.Stock{}, err
	}

	stock = models.Stock{stockID, collector.Sell, collector.Rate, collector.Buy, collector.Ticker,
		collector.Account, collector.User, collector.Time}

	return stock, nil
}
