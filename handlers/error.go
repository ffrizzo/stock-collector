package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//ErrorHandler get panic erros and send JSON
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Fatal(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err})
			}
		}()

		c.Next()
	}
}
