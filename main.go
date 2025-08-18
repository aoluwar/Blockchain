package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	Timestamp int64
}

// NewTransaction creates a new transaction
func NewTransaction(sender, receiver string, amount float64) Transaction {
	return Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}
}

// CalculateHash returns the hash of a transaction
func (tx *Transaction) CalculateHash() string {
	record := fmt.Sprintf("%s%s%f%d", tx.Sender, tx.Receiver, tx.Amount, tx.Timestamp)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Block represents a block in the blockchain
type Block struct {
	Index        int64
	Timestamp    int64
	Transactions []Transaction
	PreviousHash string
	Hash         string
	Nonce        int64
}

// NewBlock creates a new block
func NewBlock(index int64, transactions []Transaction, previousHash string) Block {
	block := Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PreviousHash: previousHash,
		Nonce:        0,
	}
	block.Hash = block.CalculateHash()
	return block
}

// CalculateHash returns the hash of a block
func (b *Block) CalculateHash() string {
	record := fmt.Sprintf("%d%d%s%d", b.Index, b.Timestamp, b.PreviousHash, b.Nonce)
	for _, tx := range b.Transactions {
		record += tx.CalculateHash()
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// MineBlock mines a block with the given difficulty
func (b *Block) MineBlock(difficulty int) {
	target := ""
	for i := 0; i < difficulty; i++ {
		target += "0"
	}

	for b.Hash[:difficulty] != target {
		b.Nonce++
		b.Hash = b.CalculateHash()
	}

	fmt.Printf("Block mined: %s\n", b.Hash)
}

// Blockchain represents the blockchain
type Blockchain struct {
	Chain               []Block
	PendingTransactions []Transaction
	Difficulty          int
	MiningReward        float64
	mu                  sync.Mutex
}

// NewBlockchain creates a new blockchain
func NewBlockchain(difficulty int, miningReward float64) *Blockchain {
	blockchain := &Blockchain{
		Chain:               []Block{},
		PendingTransactions: []Transaction{},
		Difficulty:         difficulty,
		MiningReward:       miningReward,
	}
	blockchain.CreateGenesisBlock()
	return blockchain
}

// CreateGenesisBlock creates the genesis block
func (bc *Blockchain) CreateGenesisBlock() {
	genesisBlock := NewBlock(0, []Transaction{}, "0")
	bc.Chain = append(bc.Chain, genesisBlock)
}

// GetLatestBlock returns the latest block in the chain
func (bc *Blockchain) GetLatestBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

// AddTransaction adds a transaction to the pending transactions
func (bc *Blockchain) AddTransaction(transaction Transaction) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.PendingTransactions = append(bc.PendingTransactions, transaction)
}

// MinePendingTransactions mines the pending transactions
func (bc *Blockchain) MinePendingTransactions(miningRewardAddress string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Create a reward transaction for the miner
	rewardTx := NewTransaction("BLOCKCHAIN", miningRewardAddress, bc.MiningReward)
	transactionsToMine := append(bc.PendingTransactions, rewardTx)

	// Create a new block
	newBlock := NewBlock(
		int64(len(bc.Chain)),
		transactionsToMine,
		bc.GetLatestBlock().Hash,
	)

	// Mine the block
	newBlock.MineBlock(bc.Difficulty)

	// Add the block to the chain
	bc.Chain = append(bc.Chain, newBlock)

	// Clear the pending transactions
	bc.PendingTransactions = []Transaction{}
}

// GetBalance returns the balance of an address
func (bc *Blockchain) GetBalance(address string) float64 {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	balance := 0.0

	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if tx.Sender == address {
				balance -= tx.Amount
			}
			if tx.Receiver == address {
				balance += tx.Amount
			}
		}
	}

	return balance
}

// IsChainValid checks if the chain is valid
func (bc *Blockchain) IsChainValid() bool {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		previousBlock := bc.Chain[i-1]

		// Verify the hash of the current block
		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		// Verify the previous hash reference
		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}

	return true
}

// TransactionProcessor processes transactions concurrently
type TransactionProcessor struct {
	Blockchain *Blockchain
	TxPool     []Transaction
	TxChan     chan Transaction
	mu         sync.Mutex
}

// NewTransactionProcessor creates a new transaction processor
func NewTransactionProcessor(blockchain *Blockchain) *TransactionProcessor {
	return &TransactionProcessor{
		Blockchain: blockchain,
		TxPool:     []Transaction{},
		TxChan:     make(chan Transaction, 100),
		mu:         sync.Mutex{},
	}
}

// AddTransaction adds a transaction to the pool
func (tp *TransactionProcessor) AddTransaction(transaction Transaction) {
	tp.TxChan <- transaction
}

// ProcessTransactions processes transactions from the channel
func (tp *TransactionProcessor) ProcessTransactions() {
	go func() {
		for {
			tx := <-tp.TxChan

			tp.mu.Lock()
			tp.TxPool = append(tp.TxPool, tx)
			tp.mu.Unlock()

			// If we have enough transactions, process them
			if len(tp.TxPool) >= 5 {
				tp.FlushTransactions()
			}
		}
	}()

	// Periodically flush transactions even if we don't have enough
	go func() {
		for {
			time.Sleep(5 * time.Second)
			tp.FlushTransactions()
		}
	}()
}

// FlushTransactions flushes transactions to the blockchain
func (tp *TransactionProcessor) FlushTransactions() {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if len(tp.TxPool) == 0 {
		return
	}

	// Add all transactions to the blockchain
	for _, tx := range tp.TxPool {
		tp.Blockchain.AddTransaction(tx)
	}

	// Clear the pool
	tp.TxPool = []Transaction{}
}

// Miner mines blocks
type Miner struct {
	Blockchain     *Blockchain
	MiningAddress  string
	StopChan       chan struct{}
	IsMining       bool
	mu             sync.Mutex
}

// NewMiner creates a new miner
func NewMiner(blockchain *Blockchain, miningAddress string) *Miner {
	return &Miner{
		Blockchain:    blockchain,
		MiningAddress: miningAddress,
		StopChan:      make(chan struct{}),
		IsMining:      false,
		mu:            sync.Mutex{},
	}
}

// StartMining starts mining
func (m *Miner) StartMining() {
	m.mu.Lock()
	if m.IsMining {
		m.mu.Unlock()
		return
	}
	m.IsMining = true
	m.mu.Unlock()

	go func() {
		for {
			select {
			case <-m.StopChan:
				m.mu.Lock()
				m.IsMining = false
				m.mu.Unlock()
				return
			default:
				// Check if there are pending transactions
				m.Blockchain.mu.Lock()
				pendingCount := len(m.Blockchain.PendingTransactions)
				m.Blockchain.mu.Unlock()

				if pendingCount > 0 {
					fmt.Printf("Mining block with %d transactions\n", pendingCount)
					m.Blockchain.MinePendingTransactions(m.MiningAddress)
				} else {
					time.Sleep(3 * time.Second)
				}
			}
		}
	}()
}

// StopMining stops mining
func (m *Miner) StopMining() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.IsMining {
		m.StopChan <- struct{}{}
	}
}

// BlockchainNetwork represents a blockchain network
type BlockchainNetwork struct {
	Blockchain          *Blockchain
	TransactionProcessor *TransactionProcessor
	Miners              map[string]*Miner
	mu                  sync.Mutex
}

// NewBlockchainNetwork creates a new blockchain network
func NewBlockchainNetwork(difficulty int, miningReward float64) *BlockchainNetwork {
	blockchain := NewBlockchain(difficulty, miningReward)
	transactionProcessor := NewTransactionProcessor(blockchain)

	return &BlockchainNetwork{
		Blockchain:          blockchain,
		TransactionProcessor: transactionProcessor,
		Miners:              make(map[string]*Miner),
		mu:                  sync.Mutex{},
	}
}

// Start starts the network
func (bn *BlockchainNetwork) Start() {
	bn.TransactionProcessor.ProcessTransactions()

	bn.mu.Lock()
	defer bn.mu.Unlock()

	for _, miner := range bn.Miners {
		miner.StartMining()
	}
}

// AddMiner adds a miner to the network
func (bn *BlockchainNetwork) AddMiner(id, address string) {
	bn.mu.Lock()
	defer bn.mu.Unlock()

	miner := NewMiner(bn.Blockchain, address)
	bn.Miners[id] = miner
}

// SubmitTransaction submits a transaction to the network
func (bn *BlockchainNetwork) SubmitTransaction(sender, receiver string, amount float64) {
	transaction := NewTransaction(sender, receiver, amount)
	bn.TransactionProcessor.AddTransaction(transaction)
}

// GetBalance returns the balance of an address
func (bn *BlockchainNetwork) GetBalance(address string) float64 {
	return bn.Blockchain.GetBalance(address)
}

// PrintChain prints the blockchain
func (bn *BlockchainNetwork) PrintChain() {
	fmt.Printf("Blockchain with %d blocks\n", len(bn.Blockchain.Chain))

	for i, block := range bn.Blockchain.Chain {
		// Safely get hash prefixes
		hashPrefix := block.Hash
		if len(hashPrefix) > 10 {
			hashPrefix = hashPrefix[:10]
		}
		
		prevHashPrefix := block.PreviousHash
		if len(prevHashPrefix) > 10 {
			prevHashPrefix = prevHashPrefix[:10]
		}
		
		fmt.Printf("Block %d: Index=%d, Hash=%s, PrevHash=%s, Nonce=%d, Transactions=%d\n",
			i, block.Index, hashPrefix, prevHashPrefix, block.Nonce, len(block.Transactions))
	}
}

func main() {
	// Create a new blockchain network with difficulty 2 and mining reward 100
	network := NewBlockchainNetwork(2, 100.0)

	// Add miners
	network.AddMiner("miner1", "miner1_address")
	network.AddMiner("miner2", "miner2_address")

	// Start the network
	network.Start()

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Spawn a goroutine to simulate transactions
	go func() {
		defer wg.Done()

		// Submit some transactions
		network.SubmitTransaction("alice", "bob", 50.0)
		network.SubmitTransaction("bob", "charlie", 25.0)
		network.SubmitTransaction("charlie", "alice", 10.0)

		// Wait for mining to happen
		time.Sleep(10 * time.Second)

		// Submit more transactions
		network.SubmitTransaction("alice", "charlie", 15.0)
		network.SubmitTransaction("bob", "alice", 5.0)

		// Wait for mining to happen
		time.Sleep(10 * time.Second)

		// Print the chain
		network.PrintChain()

		// Print balances
		fmt.Printf("Alice's balance: %f\n", network.GetBalance("alice"))
		fmt.Printf("Bob's balance: %f\n", network.GetBalance("bob"))
		fmt.Printf("Charlie's balance: %f\n", network.GetBalance("charlie"))
		fmt.Printf("Miner1's balance: %f\n", network.GetBalance("miner1_address"))
		fmt.Printf("Miner2's balance: %f\n", network.GetBalance("miner2_address"))
	}()

	// Wait for the simulation to complete
	wg.Wait()
	fmt.Println("Blockchain simulation completed")
}