package main

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	empHandler "mekari-test/handler"
	"mekari-test/repository"
	"mekari-test/service"
	"net/http"
	"time"
)

//go:embed migrations
var migrations embed.FS

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logrus.Errorf(err.Error())
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// Setup Logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	dbHost := viper.GetString(`database.host`)
	dbUser := viper.GetString(`database.user`)
	dbName := viper.GetString(`database.name`)
	dbPort := viper.GetString(`database.port`)
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbName, dbPort)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	logrus.Info("Pong from db")

	goose.SetBaseFS(migrations)

	if err = goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err = goose.Up(db.DB, "migrations"); err != nil {
		panic(err)
	}

	timeoutCtx := viper.GetInt(`context.timeout`)
	employeeRepo := repository.NewEmployeeRepository(db)
	employeeService := service.NewEmployeeService(employeeRepo)
	employeeHandler := empHandler.NewEmployeeHandler(employeeService, time.Duration(timeoutCtx)*time.Second)

	serverPort := viper.GetString(`server.address`)
	handler := httprouter.New()
	handler.GET("/employees", employeeHandler.GetAllEmployee)
	handler.GET("/employees/:id", employeeHandler.GetEmployeeById)
	handler.POST("/employees", employeeHandler.Register)
	handler.PUT("/employees/:id", employeeHandler.Update)
	handler.DELETE("/employees/:id", employeeHandler.Delete)

	logrus.Infof("Server run on localhost%v", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, handler))
}
