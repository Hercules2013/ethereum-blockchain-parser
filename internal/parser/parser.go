package parser

import (
	"bytes"
	"encoding/json"
	"ethereum-parser/shared"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Aliases for shared types
type (
	Transaction    = shared.Transaction
	JSONRPCRequest = shared.JSONRPCRequest
	Block          = shared.Block
)

// Parser manages the state and configuration for transaction parsing
type Parser struct {
	currentBlock        int64
	subscribed          map[string][]Transaction
	mu                  sync.RWMutex
	httpClient          *http.Client
	rpcURL              string
	processedTransactions map[string]map[string]struct{}
}

// NewParser initializes a new Parser instance and starts block scanning
func NewParser(cfg shared.Config) *Parser {
	p := &Parser{
		currentBlock:        0,
		subscribed:          make(map[string][]Transaction),
		httpClient:          &http.Client{Timeout: 10 * time.Second},
		rpcURL:              cfg.RPCURL,
		processedTransactions: make(map[string]map[string]struct{}),
	}

	go p.scanBlocks()

	return p
}

// GetCurrentBlock returns the latest block number in hexadecimal format
func (p *Parser) GetCurrentBlock() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	blockNum, err := p.fetchCurrentBlock()
	if err != nil {
		log.Println("Error fetching current block in GetCurrentBlock:", err)
	}

	p.currentBlock = blockNum
	return shared.CurrentBlockToHex(p.currentBlock)
}

// Subscribe adds an address to the list of subscribed addresses for transaction notifications
func (p *Parser) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.subscribed[address]; !exists {
		p.subscribed[address] = []Transaction{}
		return true
	}
	return false
}

// GetTransactions retrieves the list of transactions for a subscribed address
func (p *Parser) GetTransactions(address string) []Transaction {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.subscribed[address]
}

// fetchCurrentBlock gets the latest block number from the Ethereum JSON-RPC endpoint
func (p *Parser) fetchCurrentBlock() (int64, error) {
	type Response struct {
		Result string `json:"result"`
	}

	reqBody := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}`
	req, err := http.NewRequest("POST", p.rpcURL, strings.NewReader(reqBody))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}

	blockNumber, err := strconv.ParseInt(r.Result, 0, 64)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

// scanBlocks continuously fetches and processes new blocks for transactions
func (p *Parser) scanBlocks() {
	for {
		p.mu.Lock()
		startBlock := p.currentBlock
		p.mu.Unlock()

		lastBlockNumber, err := p.fetchCurrentBlock()
		if err != nil {
			log.Println("Error fetching current block in scanBlocks:", err)
			continue
		}

		if startBlock == 0 {
			continue
		}

		for i := startBlock; i <= lastBlockNumber; i++ {
			blockTransactions, err := p.getBlockTransactions(i)
			if err != nil {
				continue
			}
			p.processBlockTransactions(blockTransactions)
		}

		time.Sleep(10 * time.Second)
	}
}

// processBlockTransactions processes the transactions in a block and updates the subscriptions
func (p *Parser) processBlockTransactions(transactions []Transaction) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, tx := range transactions {
		txKey := tx.BlockNumber + "_" + tx.TransactionIndex

		p.processTransaction(tx, tx.From, txKey)
		p.processTransaction(tx, tx.To, txKey)
	}
}

// processTransaction processes a single transaction for a given address and transaction key
func (p *Parser) processTransaction(tx Transaction, address, txKey string) {
	if _, exists := p.subscribed[address]; exists {
		if _, processed := p.processedTransactions[address][txKey]; !exists && !processed {
			p.subscribed[address] = append(p.subscribed[address], tx)

			if p.processedTransactions[address] == nil {
				p.processedTransactions[address] = make(map[string]struct{})
			}
			p.processedTransactions[address][txKey] = struct{}{}
		}
	}
}

// getBlockTransactions fetches the transactions for a given block number
func (p *Parser) getBlockTransactions(blockNumber int64) ([]Transaction, error) {
	type Response struct {
		Result json.RawMessage `json:"result"`
	}

	reqBody, _ := json.Marshal(JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{shared.CurrentBlockToHex(blockNumber), true},
		ID:      83,
	})

	req, err := http.NewRequest("POST", p.rpcURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	var block Block
	if err := json.Unmarshal(r.Result, &block); err != nil {
		return nil, err
	}

	return block.Transactions, nil
}
