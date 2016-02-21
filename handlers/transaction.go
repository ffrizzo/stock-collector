package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

//TransactionHandler manages transactions for POST/DELETE/PUT
func TransactionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			tx, err := db.Begin()
			if err != nil {
				panic(err)
			}
			c.Set("tx", tx)

			defer func() {
				if err := recover(); err != nil {
					tx.Rollback()
					panic(err)
				} else {
					tx.Commit()
				}

			}()
		}

		c.Next()
	}
}
