package api

import (
	"encoding/json"
	"ethereum-parser/config"
	"ethereum-parser/internal/parser"
	"log"
	"net/http"
)

// Global parser instance
var p *parser.Parser

// init initializes the parser with configuration settings
func init() {
	cfg := config.LoadConfig()
	p = parser.NewParser(cfg)
}

// GetCurrentBlock handles HTTP requests to retrieve the current Ethereum block number
func GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := p.GetCurrentBlock()
	if err := json.NewEncoder(w).Encode(map[string]string{"current_block": block}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Subscribe handles HTTP requests to subscribe an address for transaction notifications
func Subscribe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address string `json:"address"`
	}

	// Decode the request body to get the address
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Subscribe the address and respond with the result
	success := p.Subscribe(req.Address)
	if err := json.NewEncoder(w).Encode(map[string]bool{"subscribed": success}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetTransactions handles HTTP requests to retrieve transactions for a subscribed address
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("GetTransactions() - Address: %s", address)

	// Retrieve transactions for the subscribed address and respond
	transactions := p.GetTransactions(address)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
