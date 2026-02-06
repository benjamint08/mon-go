# mon-go

barebones MongoDB client wrapper for go (driver v2)

[full example](https://github.com/benjamint08/mon-go/blob/main/example/example.go)

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

## CRUD Examples

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/benjamint08/mon-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID        string    `bson:"_id,omitempty"`
	Email     string    `bson:"email"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
}

func main() {
	ctx := context.Background()

	client, err := mon_go.New(mon_go.Config{
		URI:      "mongodb://localhost:27017",
		Database: "appdb",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	users := client.Collection("users")

	// Create
	_, err = users.InsertOne(ctx, User{
		Email:     "a@example.com",
		Name:      "Ada",
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Read (one)
	var found User
	if err := users.FindOne(ctx, bson.M{"email": "a@example.com"}).Decode(&found); err != nil {
		log.Fatal(err)
	}

	// Read (many)
	cur, err := users.Find(ctx, bson.M{}, options.Find().SetLimit(10))
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var u User
		if err := cur.Decode(&u); err != nil {
			log.Fatal(err)
		}
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Update
	_, err = users.UpdateOne(
		ctx,
		bson.M{"email": "a@example.com"},
		bson.M{"$set": bson.M{"name": "Ada Lovelace"}},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Delete
	_, err = users.DeleteOne(ctx, bson.M{"email": "a@example.com"})
	if err != nil {
		log.Fatal(err)
	}
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
