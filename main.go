package main

import (
	"crud-api-task/database"
	"crud-api-task/handlers"
	"crud-api-task/models"
	"crud-api-task/repositories"
	"crud-api-task/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string `mapstructure:"PORT"`
	DBConn 	string `mapstructure:"DB_CONN"`
}

var categories = []models.Category{
	{ID: 1, Name: "Electronics", Description: "Electronic devices"},
	{ID: 2, Name: "Mobile", Description: "Mobile phones and accessories"},
}


func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port: viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Setup routes
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/category", categoryHandler.HandleCategory)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Starting server on :" + config.Port)

	err = http.ListenAndServe(":" + config.Port, nil)

	if err != nil {
		fmt.Println("Server failed to start:", err)
	}	
}