use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::thread;
use std::time::{SystemTime, UNIX_EPOCH};
use std::fmt;
use std::hash::{Hash, Hasher};
use std::collections::hash_map::DefaultHasher;
use std::sync::mpsc;

// Define the structure of a transaction
#[derive(Clone, Debug)]
struct Transaction {
    sender: String,
    receiver: String,
    amount: f64,
    timestamp: u64,
}

impl Transaction {
    fn new(sender: String, receiver: String, amount: f64) -> Self {
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();
        
        Transaction {
            sender,
            receiver,
            amount,
            timestamp,
        }
    }
    
    fn calculate_hash(&self) -> String {
        let mut hasher = DefaultHasher::new();
        self.sender.hash(&mut hasher);
        self.receiver.hash(&mut hasher);
        (self.amount as u64).hash(&mut hasher);
        self.timestamp.hash(&mut hasher);
        format!("{:x}", hasher.finish())
    }
}

// Define the structure of a block
#[derive(Clone, Debug)]
struct Block {
    index: u64,
    timestamp: u64,
    transactions: Vec<Transaction>,
    previous_hash: String,
    hash: String,
    nonce: u64,
}

impl Block {
    fn new(index: u64, transactions: Vec<Transaction>, previous_hash: String) -> Self {
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();
        
        let mut block = Block {
            index,
            timestamp,
            transactions,
            previous_hash,
            hash: String::new(),
            nonce: 0,
        };
        
        block.hash = block.calculate_hash();
        block
    }
    
    fn calculate_hash(&self) -> String {
        let mut hasher = DefaultHasher::new();
        self.index.hash(&mut hasher);
        self.timestamp.hash(&mut hasher);
        for tx in &self.transactions {
            tx.calculate_hash().hash(&mut hasher);
        }
        self.previous_hash.hash(&mut hasher);
        self.nonce.hash(&mut hasher);
        format!("{:x}", hasher.finish())
    }
    
    fn mine_block(&mut self, difficulty: usize) {
        let target = "0".repeat(difficulty);
        
        while &self.hash[0..difficulty] != target {
            self.nonce += 1;
            self.hash = self.calculate_hash();
        }
        
        println!("Block mined: {}", self.hash);
    }
}

// Define the blockchain
#[derive(Clone)]
struct Blockchain {
    chain: Vec<Block>,
    pending_transactions: Vec<Transaction>,
    difficulty: usize,
    mining_reward: f64,
}

impl Blockchain {
    fn new(difficulty: usize, mining_reward: f64) -> Self {
        let mut blockchain = Blockchain {
            chain: Vec::new(),
            pending_transactions: Vec::new(),
            difficulty,
            mining_reward,
        };
        
        // Create the genesis block
        blockchain.create_genesis_block();
        blockchain
    }
    
    fn create_genesis_block(&mut self) {
        let genesis_block = Block::new(0, Vec::new(), String::from("0"));
        self.chain.push(genesis_block);
    }
    
    fn get_latest_block(&self) -> &Block {
        &self.chain[self.chain.len() - 1]
    }
    
    fn add_transaction(&mut self, transaction: Transaction) {
        self.pending_transactions.push(transaction);
    }
    
    fn mine_pending_transactions(&mut self, mining_reward_address: String) {
        // Create a reward transaction for the miner
        let reward_tx = Transaction::new(
            String::from("BLOCKCHAIN"),
            mining_reward_address,
            self.mining_reward,
        );
        
        let mut transactions_to_mine = self.pending_transactions.clone();
        transactions_to_mine.push(reward_tx);
        
        let new_block = Block::new(
            self.chain.len() as u64,
            transactions_to_mine,
            self.get_latest_block().hash.clone(),
        );
        
        let mut mineable_block = new_block.clone();
        mineable_block.mine_block(self.difficulty);
        
        self.chain.push(mineable_block);
        self.pending_transactions = Vec::new();
    }
    
    fn get_balance(&self, address: &str) -> f64 {
        let mut balance = 0.0;
        
        for block in &self.chain {
            for tx in &block.transactions {
                if tx.sender == address {
                    balance -= tx.amount;
                }
                
                if tx.receiver == address {
                    balance += tx.amount;
                }
            }
        }
        
        balance
    }
    
    fn is_chain_valid(&self) -> bool {
        for i in 1..self.chain.len() {
            let current_block = &self.chain[i];
            let previous_block = &self.chain[i - 1];
            
            // Verify the hash of the current block
            if current_block.hash != current_block.calculate_hash() {
                return false;
            }
            
            // Verify the previous hash reference
            if current_block.previous_hash != previous_block.hash {
                return false;
            }
        }
        
        true
    }
}

impl fmt::Debug for Blockchain {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "Blockchain with {} blocks", self.chain.len())
    }
}

// Concurrent transaction processor
struct TransactionProcessor {
    blockchain: Arc<Mutex<Blockchain>>,
    tx_pool: Arc<Mutex<Vec<Transaction>>>,
}

impl TransactionProcessor {
    fn new(blockchain: Arc<Mutex<Blockchain>>) -> Self {
        TransactionProcessor {
            blockchain,
            tx_pool: Arc::new(Mutex::new(Vec::new())),
        }
    }
    
    fn add_transaction(&self, transaction: Transaction) {
        let mut tx_pool = self.tx_pool.lock().unwrap();
        tx_pool.push(transaction);
    }
    
    fn process_transactions(&self) {
        let tx_pool_clone = self.tx_pool.clone();
        let blockchain_clone = self.blockchain.clone();
        
        thread::spawn(move || {
            loop {
                let transactions = {
                    let mut tx_pool = tx_pool_clone.lock().unwrap();
                    if tx_pool.is_empty() {
                        thread::sleep(std::time::Duration::from_secs(1));
                        continue;
                    }
                    
                    let txs = tx_pool.clone();
                    tx_pool.clear();
                    txs
                };
                
                let mut blockchain = blockchain_clone.lock().unwrap();
                for tx in transactions {
                    blockchain.add_transaction(tx);
                }
            }
        });
    }
}

// Concurrent miner
struct Miner {
    blockchain: Arc<Mutex<Blockchain>>,
    mining_address: String,
    stop_signal: Arc<Mutex<bool>>,
}

impl Miner {
    fn new(blockchain: Arc<Mutex<Blockchain>>, mining_address: String) -> Self {
        Miner {
            blockchain,
            mining_address,
            stop_signal: Arc::new(Mutex::new(false)),
        }
    }
    
    fn start_mining(&self) -> thread::JoinHandle<()> {
        let blockchain_clone = self.blockchain.clone();
        let mining_address = self.mining_address.clone();
        let stop_signal = self.stop_signal.clone();
        
        thread::spawn(move || {
            loop {
                {
                    let should_stop = *stop_signal.lock().unwrap();
                    if should_stop {
                        break;
                    }
                }
                
                let mut blockchain = blockchain_clone.lock().unwrap();
                let pending_count = blockchain.pending_transactions.len();
                
                if pending_count > 0 {
                    println!("Mining block with {} transactions", pending_count);
                    blockchain.mine_pending_transactions(mining_address.clone());
                } else {
                    drop(blockchain); // Release the lock
                    thread::sleep(std::time::Duration::from_secs(3));
                }
            }
        })
    }
    
    fn stop_mining(&self) {
        let mut stop = self.stop_signal.lock().unwrap();
        *stop = true;
    }
}

// Blockchain network simulator
struct BlockchainNetwork {
    blockchain: Arc<Mutex<Blockchain>>,
    transaction_processor: TransactionProcessor,
    miners: HashMap<String, Miner>,
}

impl BlockchainNetwork {
    fn new(difficulty: usize, mining_reward: f64) -> Self {
        let blockchain = Arc::new(Mutex::new(Blockchain::new(difficulty, mining_reward)));
        let transaction_processor = TransactionProcessor::new(blockchain.clone());
        
        BlockchainNetwork {
            blockchain,
            transaction_processor,
            miners: HashMap::new(),
        }
    }
    
    fn start(&mut self) {
        self.transaction_processor.process_transactions();
        
        for (_, miner) in &self.miners {
            miner.start_mining();
        }
    }
    
    fn add_miner(&mut self, id: String, address: String) {
        let miner = Miner::new(self.blockchain.clone(), address);
        self.miners.insert(id, miner);
    }
    
    fn submit_transaction(&self, sender: String, receiver: String, amount: f64) {
        let transaction = Transaction::new(sender, receiver, amount);
        self.transaction_processor.add_transaction(transaction);
    }
    
    fn get_balance(&self, address: &str) -> f64 {
        let blockchain = self.blockchain.lock().unwrap();
        blockchain.get_balance(address)
    }
    
    fn print_chain(&self) {
        let blockchain = self.blockchain.lock().unwrap();
        println!("Blockchain: {:?}", blockchain);
        
        for (i, block) in blockchain.chain.iter().enumerate() {
            println!("Block {}: {:?}", i, block);
        }
    }
}

fn main() {
    // Create a new blockchain network with difficulty 2 and mining reward 100
    let mut network = BlockchainNetwork::new(2, 100.0);
    
    // Add miners
    network.add_miner(String::from("miner1"), String::from("miner1_address"));
    network.add_miner(String::from("miner2"), String::from("miner2_address"));
    
    // Start the network
    network.start();
    
    // Create a channel for communication between threads
    let (tx, rx) = mpsc::channel();
    
    // Spawn a thread to simulate transactions
    let network_clone = Arc::new(Mutex::new(network));
    let tx_clone = tx.clone();
    
    thread::spawn(move || {
        let network = network_clone.lock().unwrap();
        
        // Submit some transactions
        network.submit_transaction(String::from("alice"), String::from("bob"), 50.0);
        network.submit_transaction(String::from("bob"), String::from("charlie"), 25.0);
        network.submit_transaction(String::from("charlie"), String::from("alice"), 10.0);
        
        // Wait for mining to happen
        thread::sleep(std::time::Duration::from_secs(10));
        
        // Submit more transactions
        network.submit_transaction(String::from("alice"), String::from("charlie"), 15.0);
        network.submit_transaction(String::from("bob"), String::from("alice"), 5.0);
        
        // Wait for mining to happen
        thread::sleep(std::time::Duration::from_secs(10));
        
        // Print the chain
        network.print_chain();
        
        // Print balances
        println!("Alice's balance: {}", network.get_balance("alice"));
        println!("Bob's balance: {}", network.get_balance("bob"));
        println!("Charlie's balance: {}", network.get_balance("charlie"));
        println!("Miner1's balance: {}", network.get_balance("miner1_address"));
        println!("Miner2's balance: {}", network.get_balance("miner2_address"));
        
        // Signal main thread to exit
        tx_clone.send(()).unwrap();
    });
    
    // Wait for the simulation to complete
    rx.recv().unwrap();
    println!("Blockchain simulation completed");
}