# Crypto Arbitrage Bot

A sophisticated cryptocurrency arbitrage bot that monitors mempool transactions and executes arbitrage opportunities with advanced security features.

## Features

- Real-time mempool transaction monitoring
- Proof of Work (PoW) validation for node verification
- AES encryption for secure transaction data
- Integration with Sentry for error tracking
- DataDog monitoring and tracing
- Structured logging system

## Prerequisites

- Go 1.20 or higher
- DataDog agent running locally (for metrics)
- Sentry account (for error tracking)
- Access to cryptocurrency DEX API

## Setup Instructions

1. Clone the repository:
```bash
git clone https://github.com/ivan-cody/crypto_bot.git
cd crypto_bot
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure the application:
   - Copy `config/config.yaml.example` to `config/config.yaml`
   - Update the following values:
     - `mempool_url`: Your mempool endpoint
     - `sentry_dsn`: Your Sentry DSN
     - `encryption_key`: 32-byte hex-encoded key
     - `datadog.api_key`: Your DataDog API key
     - `datadog.app_key`: Your DataDog APP key

4. Start the DataDog agent locally

5. Run the bot:
```bash
go run main.go
```

## Technical Details

### Architecture Flow

1. **Initialization**
   - Load configuration
   - Set up encryption (AES-GCM)
   - Initialize monitoring (Sentry & DataDog)
   - Connect to mempool

2. **Transaction Processing**
   - Monitor mempool for new transactions
   - Analyze each transaction for arbitrage opportunities
   - Validate opportunities against minimum amount thresholds

3. **Arbitrage Execution**
   - Perform PoW validation (4 leading zeros required)
   - Encrypt transaction data
   - Execute buy order on DEX
   - Wait for confirmation
   - Execute sell order
   - Log results

### Security Features

- **Proof of Work**: Validates node authenticity before transactions
- **Encryption**: AES-GCM encryption for sensitive data
- **Monitoring**: Real-time error tracking and performance monitoring
- **Tracing**: Distributed tracing for transaction flow

### Performance Monitoring

- Memory usage tracking
- Goroutine count monitoring
- Transaction processing metrics
- Custom DataDog dashboards

## Testing

Run the test suite:
```bash
go test ./...
```

## Production Deployment

For production deployment:
1. Set `environment: "production"` in config
2. Use secure encryption keys
3. Configure proper monitoring thresholds
4. Set up alerting in Sentry and DataDog

## License

MIT License - see LICENSE file for details 