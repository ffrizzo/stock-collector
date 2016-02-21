package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"fmt"

	"github.com/ffrizzo/stock-collector/controllers"
	"github.com/ffrizzo/stock-collector/handlers"
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
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.name"))
	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		panic(err)
	}

	session.db = db
}

func main() {
	defer session.db.Close()
	g := gin.New()

	g.Use(handlers.ErrorHandler())
	g.Use(handlers.TransactionHandler(session.db))

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
