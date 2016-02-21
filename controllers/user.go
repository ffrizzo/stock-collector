package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ffrizzo/stock-collector/utils"
	"github.com/gin-gonic/gin"
)

//UserController represents the controller for User resource
type UserController struct {
	db *sql.DB
}

//NewUserController instantiate a new controller for user
func NewUserController(db *sql.DB) *UserController {
	return &UserController{db}
}

//UserMostActive get the user most active on the last hour
func (controller UserController) UserMostActive(c *gin.Context) {

	endTime := time.Now()
	startTime := util.SubtractHours(endTime, 1)

	query := `select username, count(id) from stock where time between $1 and $2
        group by username order by count desc limit 1`
	var username string
	var count int
	err := controller.db.QueryRow(query, startTime, endTime).Scan(&username, &count)
	if err != nil && err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"success": true, "message": "Does not have activity since one hour ago."})
		return
	} else if err != nil {
		panic(err)
	}

	c.JSON(http.StatusBadRequest, gin.H{"success": true, "username": username, "total": count})
}
