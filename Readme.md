# Ethereum Parser

A lightweight and efficient Ethereum parser API written in Go.

## Features

- Retrieve the current block number
- Subscribe to an address for transaction notifications
- Retrieve transactions for a subscribed address

## Project Structure

The project is organized as follows:

```
ethereum-parser/
|-- cmd/
|   |-- main.go
|-- config/
|   |-- config.go
|-- internal/
|   |-- api/
|   |   |-- api.go
|   |-- parser/
|       |-- parser.go
|       |-- parser_test.go
|-- shared/
|   |-- types.go
|   |-- utils.go
|-- Readme.md
|-- go.mod
|-- go.sum
```

## Configuration

Configure the application using the following environment variables:

- `SERVER_PORT`: The port on which the server will run (default: `8080`)
- `RPC_URL`: The Ethereum JSON-RPC URL (default: `https://cloudflare-eth.com`)

## Running Locally

1. Clone the repository:

    ```sh
    git clone https://github.com/Hercules2013/ethereum-blockchain-parser.git
    cd ethereum-blockchain-parser
    ```

2. Run the application:

    ```sh
    go run cmd/main.go
    ```

## API Endpoints

- `GET /current_block`: Retrieve the current block number.
- `POST /subscribe`: Subscribe to an address for transaction notifications. Request body: `{ "address": "0x123" }`
- `GET /transactions?address=0x123`: Retrieve transactions for a subscribed address.

## Testing

Run the tests using:

```sh
go test ./internal/parser/ -v