package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ffrizzo/stock-collector/models"
	"github.com/ffrizzo/stock-collector/utils"
	"github.com/gin-gonic/gin"
)

//StockController represents the controller for User resource
type StockController struct {
	db *sql.DB
}

//NewStockController instantiate a new controller for user
func NewStockController(db *sql.DB) *StockController {
	return &StockController{db}
}

//StockMostExpensive get the most expensive stock for each account
func (controller StockController) StockMostExpensive(c *gin.Context) {
	startTime, endTime := util.GetStartAndEndTimeOfYesterday()
	fmt.Println(startTime)
	fmt.Println(endTime)

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
func (controller StockController) StockMeanMedian(c *gin.Context) {

	account := c.Query("account")
	ticker := c.Query("ticker")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	query := `select median(sell) sell_media, median(rate) rate_media, median(buy) buy_median,
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
		"stock": gin.H{"ticker": ticker, "account": account,
			"sell_median": sellMedia, "sell_avg": sellAvg,
			"rate_median": rateMedia, "rate_avg": rateAvg,
			"buy_median": buyMedia, "buy_avg": buyAvg},
	})
}
