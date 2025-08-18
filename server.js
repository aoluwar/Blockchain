const http = require('http');
const fs = require('fs');
const path = require('path');

const PORT = 3000;

const server = http.createServer((req, res) => {
    // Handle the root path by serving the interface.html file
    if (req.url === '/' || req.url === '/index.html') {
        const filePath = path.join(__dirname, 'interface.html');
        fs.readFile(filePath, (err, content) => {
            if (err) {
                res.writeHead(500);
                res.end(`Error loading the interface: ${err.message}`);
                return;
            }
            res.writeHead(200, { 'Content-Type': 'text/html' });
            res.end(content, 'utf-8');
        });
        return;
    }

    // For API endpoints (can be expanded later)
    if (req.url.startsWith('/api/')) {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        
        // Mock API responses
        if (req.url === '/api/blockchain') {
            res.end(JSON.stringify({ message: 'Blockchain API endpoint' }));
        } else if (req.url === '/api/transactions') {
            res.end(JSON.stringify({ message: 'Transactions API endpoint' }));
        } else {
            res.writeHead(404);
            res.end(JSON.stringify({ error: 'API endpoint not found' }));
        }
        return;
    }

    // Handle 404 for any other routes
    res.writeHead(404);
    res.end('Not Found');
});

server.listen(PORT, () => {
    console.log(`Server running at http://localhost:${PORT}/`);
    console.log(`Open your browser to http://localhost:${PORT}/ to view the blockchain interface`);
});