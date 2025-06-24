package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Item struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var db *gorm.DB

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}

	db.AutoMigrate(&Item{})

	router := mux.NewRouter()
	router.HandleFunc("/inventory", getItems).Methods("GET")
	router.HandleFunc("/inventory", createItem).Methods("POST")
	router.HandleFunc("/inventory/{id}/reduce", reduceStock).Methods("PUT")

	log.Println("Inventory service running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}

func getItems(w http.ResponseWriter, r *http.Request) {
	var items []Item
	db.Find(&items)
	json.NewEncoder(w).Encode(items)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	json.NewDecoder(r.Body).Decode(&item)
	item.ID = uuid.New()
	db.Create(&item)
	json.NewEncoder(w).Encode(item)
}

func reduceStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item Item
	if err := db.First(&item, "id = ?", id).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	item.Stock -= 1
	db.Save(&item)
	json.NewEncoder(w).Encode(item)
}
