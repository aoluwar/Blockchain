# Blockchain Concurrency Implementation

This project demonstrates a simple blockchain implementation with concurrency features using both Rust and Go. The implementations showcase how each language handles concurrent operations in a blockchain context.

## Project Structure

- `main.rs` - Rust implementation of the blockchain
- `main.go` - Go implementation of the blockchain
- `interface.html` - Web interface for interacting with the blockchain
- `server.js` - Node.js server to serve the web interface
- `start_interface.bat` - Batch script to start the web interface server

## Features

Both implementations include:

- Basic blockchain structure with blocks and transactions
- Proof-of-work mining algorithm
- Concurrent transaction processing
- Concurrent block mining
- Balance tracking for addresses
- Chain validation

## Concurrency Features

### Rust Implementation

The Rust implementation uses:

- `Arc<Mutex<T>>` for shared state between threads
- `thread::spawn` for creating worker threads
- `mpsc` channels for communication between threads
- Thread-safe data structures with explicit locking

### Go Implementation

The Go implementation uses:

- Goroutines for concurrent execution
- Channels for communication between goroutines
- `sync.Mutex` for protecting shared state
- `sync.WaitGroup` for synchronization

## How to Run

### Rust Implementation

```bash
# Navigate to the project directory
cd blockchain

# Run the Rust implementation
rustc main.rs
./main
```

### Go Implementation

```bash
# Navigate to the project directory
cd blockchain

# Run the Go implementation
go run main.go
```

### Web Interface

```bash
# Navigate to the project directory
cd blockchain

# Start the web interface server
.\start_interface.bat
```

Then open your browser and navigate to http://localhost:3000/ to access the blockchain interface.

## Key Differences

### Concurrency Model

- **Rust**: Uses a thread-based concurrency model with explicit ownership and borrowing rules enforced by the compiler. Shared state is managed through `Arc<Mutex<T>>` to ensure thread safety.

- **Go**: Uses a lightweight goroutine-based concurrency model with channels for communication. The "share memory by communicating" philosophy is emphasized, though mutexes are also available when needed.

### Memory Safety

- **Rust**: Provides memory safety guarantees at compile time through its ownership system, preventing data races and other concurrency bugs.

- **Go**: Provides memory safety through garbage collection and runtime checks, making concurrent programming more accessible but with potential runtime overhead.

### Performance

- **Rust**: Generally offers better performance due to zero-cost abstractions and fine-grained control over memory, but requires more explicit code for concurrency.

- **Go**: Offers good performance with simpler concurrency primitives, but may have higher memory usage due to garbage collection.

## Learning Points

1. How blockchain data structures can be implemented in different languages
2. Different approaches to concurrency in systems programming
3. Trade-offs between safety, performance, and simplicity in concurrent code
4. Practical implementation of proof-of-work consensus algorithm