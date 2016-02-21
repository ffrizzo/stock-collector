package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"fmt"

	"github.com/ffrizzo/stock-colector/controllers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var session struct {
	db *sql.DB
}

func init() {
	InitConfig()
	SetDefault()
	InitDB()
}

//InitConfig initialize viper
func InitConfig() {
	viper.SetConfigType("yml")
	viper.AddConfigPath("configs/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error on config file: %s\n", err))
	}
}

//SetDefault values of configurations
func SetDefault() {
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.name", "stockCollector")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
}

//InitDB start the connection with Database
func InitDB() {
	connectionURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		viper.GetString("database.user"), viper.GetString("database.password"),
		viper.GetString("database.host"), viper.GetString("database.port"),
		viper.GetString("database.name"))
	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		panic(err)
	}

	session.db = db
}

//TransactionHandler manages transactions for POST/DELETE/PUT
func TransactionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			tx, err := session.db.Begin()
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

func main() {
	defer session.db.Close()
	g := gin.New()

	g.Use(ErrorHandler())
	g.Use(TransactionHandler())

	cc := controllers.NewCollectorController(session.db)
	uc := controllers.NewUserController(session.db)
	sc := controllers.NewStockController(session.db)

	v1 := g.Group("/v1")
	{
		v1.POST("/collector", cc.Collector)
		v1.GET("/collector/user/active/last/hour", uc.UserMostActive)
		v1.GET("/collector/stock/expensive/last/day", sc.StockMostExpensive)
		v1.GET("/collector/stock/mean_media", sc.StockMeanMedian)
	}

	// port := []string{":", viper.GetString("server.port")}
	// g.Run(strings.Join(port, ""))
	g.Run()
}
