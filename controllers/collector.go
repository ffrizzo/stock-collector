package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/ffrizzo/stock-colector/models"
	"github.com/ffrizzo/stock-colector/utils"
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

	stock, err := SaveStock(controller.db, collector)
	if err != nil {
		panic(fmt.Sprintf("Error to retrieve data for stock. Error: %s", err))
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "stock": stock})
}

//SaveStock save or update data for Stock
func SaveStock(db *sql.DB, collector models.Stock) (models.Stock, error) {
	stock := models.Stock{}

	query := `insert into stock (sell, rate, buy, ticker, account, username, time)
        values ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	var stockID int64
	err := db.QueryRow(query, collector.Sell, collector.Rate, collector.Buy, collector.Ticker,
		collector.Account, collector.User, collector.Time).
		Scan(&stockID)
	if err != nil {
		return models.Stock{}, err
	}

	stock = models.Stock{stockID, collector.Sell, collector.Rate, collector.Buy, collector.Ticker,
		collector.Account, collector.User, collector.Time}

	return stock, nil
}

//UserMostActive get the most active user on last hour
func (controller CollectorController) UserMostActive(c *gin.Context) {

	startTime := time.Now()
	endTime := util.SubtractHours(startTime, 1)

	query := `select username, count(id) from stock where time between $1 and $2
        group by username order by count desc limit 1`
	var username string
	var count int
	err := controller.db.QueryRow(query, startTime, endTime).Scan(&username, &count)
	if err != nil && err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"success": true, "message": "Does not have user active on the last hour"})
		return
	} else if err != nil {
		panic(err)
	}

	c.JSON(http.StatusBadRequest, gin.H{"success": true, "username": username, "total": count})
}

//StockMostExpensive get the most expensive stock for each account
func (controller CollectorController) StockMostExpensive(c *gin.Context) {

	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	query := `select id, sell, rate, buy, ticker, account, username, time from stock
        where buy in (select max(buy) as max_buy from stock
                      where time between $1 and $2
                      group by account)`
	rows, err := controller.db.Query(query, startTime, endTime)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		stock := models.Stock{}
		rows.Scan(&stock.ID, &stock.Sell, &stock.Rate, &stock.Buy, &stock.Ticker,
			&stock.Account, &stock.User, &stock.Time)

		stocks = append(stocks, stock)
	}

	c.JSON(http.StatusBadRequest, gin.H{"success": true, "stocks": stocks})
}

//StockMeanMedian get the meand and median in an interval of dates
func (controller CollectorController) StockMeanMedian(c *gin.Context) {

	account := c.Query("account")
	ticker := c.Query("ticker")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	query := `sselect median(sell) sell_media, median(rate) rate_media, median(buy) buy_median,
              avg(sell) sell_avg, avg(rate) rate_avg, avg(buy) buy_avg from stock
              where time between $1 and $2
              and account = $3
              and ticker = $4
              group by ticker, account`

	var sellMedia, rateMedia, buyMedia float64
	var sellAvg, rateAvg, buyAvg float64

	err := controller.db.QueryRow(query, startTime, endTime, account, ticker).
		Scan(&sellMedia, &rateMedia, &buyMedia, &sellAvg, &rateAvg, &buyAvg)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusBadRequest, gin.H{"success": true,
		"stock": gin.H{"ticker": ticker,
			"sell_median": sellMedia, "sell_avg": sellAvg,
			"rate_median": rateMedia, "rate_avg": rateAvg,
			"buy_median": buyMedia, "buy_avg": buyAvg},
	})
}
