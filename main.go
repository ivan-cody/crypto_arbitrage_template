package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/getsentry/sentry-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"crypto_bot/arbitrage"
	"crypto_bot/config"
	"crypto_bot/logger"
	"crypto_bot/mempool"
)

func initMonitoring(cfg config.Config) error {
	// Initialize Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.Environment,
		Debug:            true,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		return err
	}

	// Initialize DataDog tracer
	tracer.Start(
		tracer.WithService("crypto-arbitrage-bot"),
		tracer.WithEnv(cfg.Environment),
		tracer.WithServiceVersion("1.0.0"),
	)

	// Initialize DataDog metrics
	statsdClient, err := statsd.New("localhost:8125",
		statsd.WithNamespace("crypto.arbitrage."),
		statsd.WithTags([]string{"env:" + cfg.Environment}),
	)
	if err != nil {
		return err
	}

	// Start periodic metrics reporting
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			statsdClient.Gauge("system.memory", getMemoryUsage(), nil, 1)
			statsdClient.Gauge("system.goroutines", float64(getGoroutineCount()), nil, 1)
		}
	}()

	return nil
}

func getMemoryUsage() float64 {
	// Placeholder for actual memory metrics
	return 0.0
}

func getGoroutineCount() int {
	// Placeholder for actual goroutine count
	return 0
}

func initEncryption(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm, nil
}

func main() {
	// Initialize logging
	logger.Init("bot.log")
	log.Println("Starting Crypto Arbitrage Bot")

	defer sentry.Flush(2 * time.Second)
	defer tracer.Stop()

	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize monitoring
	if err := initMonitoring(cfg); err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Failed to initialize monitoring: %v", err)
	}

	// Initialize encryption
	encryptionKey, err := hex.DecodeString(cfg.EncryptionKey)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Invalid encryption key: %v", err)
	}

	cipher, err := initEncryption(encryptionKey)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Failed to initialize encryption: %v", err)
	}

	// Connect to mempool
	mempoolConn, err := mempool.Connect(cfg.MempoolURL)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Failed to connect to mempool: %v", err)
	}
	defer mempoolConn.Close()

	// Start transaction processing
	arbitrageProcessor := arbitrage.NewProcessor(cfg, cipher)
	span := tracer.StartSpan("process.transactions")
	defer span.Finish()

	for tx := range mempoolConn.GetTransactions() {
		go func(tx mempool.Transaction) {
			txSpan := tracer.StartSpan("process.single.transaction", tracer.ChildOf(span.Context()))
			defer txSpan.Finish()
			arbitrageProcessor.ProcessTransaction(tx)
		}(tx)
	}
}
