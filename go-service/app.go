package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
)

const stateStoreName = "statestore"

type Price struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type InventoryItem struct {
	ID       string  `json:"id"`
	Item     string  `json:"item"`
	Location string  `json:"location"`
	Priority string  `json:"priority"`
	Price    *Price  `json:"price,omitempty"`
}

type App struct {
	Router     *mux.Router
	daprClient dapr.Client
}

func (a *App) Initialize(client dapr.Client) {
	a.daprClient = client
	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/", a.Hello).Methods("GET")
	a.Router.HandleFunc("/inventory", a.GetInventory).Methods("GET")
	a.Router.HandleFunc("/inventory", a.CreateOrUpdateInventory).Methods("POST", "PUT")
}

func (a *App) Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world! It's me"))
}

func (a *App) GetInventory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	item, err := a.daprClient.GetState(context.Background(), stateStoreName, id)
	if err != nil {
		http.Error(w, "Failed to retrieve inventory", http.StatusInternalServerError)
		return
	}

	if item.Value == nil || len(item.Value) == 0 {
		http.Error(w, "Inventory item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(item.Value)
}

func (a *App) CreateOrUpdateInventory(w http.ResponseWriter, r *http.Request) {
	var item InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if item.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	if err := validateInventoryItem(&item); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		http.Error(w, "Failed to serialize inventory item", http.StatusInternalServerError)
		return
	}

	if err := a.daprClient.SaveState(context.Background(), stateStoreName, item.ID, data); err != nil {
		http.Error(w, "Failed to save inventory", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func validateInventoryItem(item *InventoryItem) error {
	if item.Price != nil {
		if item.Price.Value < 0 {
			return &ValidationError{"price.value must be non-negative"}
		}
		matched, _ := regexp.MatchString("^[A-Z]{3}$", item.Price.Currency)
		if !matched {
			return &ValidationError{"price.currency must be a 3-letter uppercase code (e.g., USD, EUR)"}
		}
	}
	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
