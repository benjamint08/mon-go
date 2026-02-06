package mon_go

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	defaultConnectTimeout = 10 * time.Second
	defaultPingTimeout    = 2 * time.Second
)

// Config holds the configuration for the MongoDB client
type Config struct {
	// connection string
	URI string
	// default database name
	Database string
	// connectTimeout bounds the initial connection attempt
	ConnectTimeout time.Duration
	// pingTimeout bounds the initial ping call
	PingTimeout time.Duration
	// maxPoolSize sets the driver's max pool size, use 0 to keep driver defaults
	MaxPoolSize uint64
}

// Client is a wrapper
type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

// New creates and verifies a MongoDB client
func New(cfg Config) (*Client, error) {
	if cfg.URI == "" {
		return nil, errors.New("mongo: missing URI")
	}
	if cfg.Database == "" {
		return nil, errors.New("mongo: missing database name")
	}

	connectTimeout := cfg.ConnectTimeout
	if connectTimeout <= 0 {
		connectTimeout = defaultConnectTimeout
	}

	pingTimeout := cfg.PingTimeout
	if pingTimeout <= 0 {
		pingTimeout = defaultPingTimeout
	}

	opts := options.Client().ApplyURI(cfg.URI)
	if cfg.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(cfg.MaxPoolSize)
	}
	if connectTimeout > 0 {
		opts.SetConnectTimeout(connectTimeout)
	}

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), pingTimeout)
	defer pingCancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}

	return &Client{
		client: client,
		db:     client.Database(cfg.Database),
	}, nil
}

// NewFromEnv builds a client from env vars
// env vars:
// - MONGO_URI: connection string
// - MONGO_DB: default database name
func NewFromEnv() (*Client, error) {
	return New(Config{
		URI:      os.Getenv("MONGO_URI"),
		Database: os.Getenv("MONGO_DB"),
	})
}

// DB returns the default database
func (c *Client) DB() *mongo.Database {
	return c.db
}

// Collection returns a collection from the default database
func (c *Client) Collection(name string) *mongo.Collection {
	return c.db.Collection(name)
}

// Raw returns the underlying mongo.Client
func (c *Client) Raw() *mongo.Client {
	return c.client
}

// Disconnect closes the client
func (c *Client) Disconnect(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}
