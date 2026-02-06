# mon-go

barebones MongoDB client wrapper for go (driver v2)

## Install

```bash
go get github.com/benjamint08/mon-go
```

## Usage

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/benjamint08/mon-go"
)

func main() {
	client, err := mon_go.New(mon_go.Config{
		URI:            "mongodb://localhost:27017",
		Database:       "appdb",
		ConnectTimeout: 5 * time.Second,
		PingTimeout:    2 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	users := client.Collection("users")
	_ = users
}
```

## Env Setup

`NewFromEnv` needs:
- `MONGO_URI`
- `MONGO_DB`

```go
package main

import (
	"context"
	"log"

	"github.com/benjamint08/mon-go"
)

func main() {
	client, err := mon_go.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
}
```

## Config Fields

- `URI`: MongoDB connection string (required)
- `Database`: default database name (required)
- `ConnectTimeout`: optional connect timeout
- `PingTimeout`: optional ping timeout
- `MaxPoolSize`: optional driver max pool size (0 keeps defaults)
