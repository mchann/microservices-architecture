package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Order struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID  uuid.UUID `json:"product_id"`
	Total      float64   `json:"total"`
	IsPaid     bool      `json:"is_paid"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PaymentResponse struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ID         uuid.UUID `json:"id"`
	Value      float64   `json:"value"`
	MerchantID uuid.UUID `json:"merchant_id"`
	Region     string    `json:"region"`
}

func main() {
	// Koneksi DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed connect DB:", err)
	}
	db.AutoMigrate(&Order{})

	// Routing
	router := mux.NewRouter()
	router.HandleFunc("/order", createOrder).Methods("POST")
	router.HandleFunc("/order/{id}/paid", markAsPaid).Methods("PUT")
	router.HandleFunc("/order", getOrders).Methods("GET")

	log.Println("Order service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order.ID = uuid.New()
	order.IsPaid = false
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// === Call Payment Service ===
	basePaymentURL := fmt.Sprintf("http://%s", os.Getenv("PAYMENT_SERVICE_HOST"))
	paymentPayload := map[string]interface{}{
		"value":       order.Total,
		"merchant_id": order.ProductID, // sementara pakai product_id sebagai merchant_id
	}
	payloadBytes, err := json.Marshal(paymentPayload)
	if err != nil {
		http.Error(w, "Failed marshal payment payload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := http.Post(basePaymentURL+"/payments/id/api/v1/", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Failed to call payment service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		log.Printf("Payment service error: %s", body)
		http.Error(w, "Payment service failed", http.StatusInternalServerError)
		return
	}

	var p PaymentResponse
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		http.Error(w, "Failed decode payment response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// === Call Inventory Service ===
	baseInventoryURL := fmt.Sprintf("http://%s", os.Getenv("INVENTORY_SERVICE_HOST"))
	req, err := http.NewRequest("PUT", baseInventoryURL+"/inventory/"+order.ProductID.String()+"/reduce", nil)
	if err != nil {
		http.Error(w, "Failed build inventory request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	invResp, err := client.Do(req)
	if err != nil || invResp.StatusCode != http.StatusOK {
		http.Error(w, "Failed reduce stock", http.StatusInternalServerError)
		return
	}
	defer invResp.Body.Close()

	// === Simpan ke DB ===
	if err := db.Create(&order).Error; err != nil {
		http.Error(w, "Failed save order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func markAsPaid(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var order Order
	if err := db.First(&order, "id = ?", id).Error; err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	order.IsPaid = true
	order.UpdatedAt = time.Now()
	db.Save(&order)
	json.NewEncoder(w).Encode(order)
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	var orders []Order
	db.Find(&orders)

	var result []string
	for _, o := range orders {
		if o.IsPaid {
			msg := fmt.Sprintf("merchant %s received payment %.6f USD", o.ProductID.String(), o.Total)
			result = append(result, msg)
		}
	}
	json.NewEncoder(w).Encode(result)
}
