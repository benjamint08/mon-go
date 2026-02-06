// package main

import (
	"context"
	"time"

	mon_go "github.com/benjamint08/mon-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func main() {
	start := time.Now()
	client, err := mon_go.New(mon_go.Config{
		URI:            "mongodb://localhost:27017",
		Database:       "mon_go_example",
		ConnectTimeout: 5 * time.Second,
		PingTimeout:    2 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	println("Connected to MongoDB in", time.Since(start).String())

	defer client.Disconnect(context.Background())

	users := client.Collection("users")

	start = time.Now()
	// Insert a user
	_, err = users.InsertOne(context.Background(), User{Name: "Alice", Age: 30})
	if err != nil {
		panic(err)
	}
	println("Inserted user in", time.Since(start).String())

	start = time.Now()
	var foundUser User
	if err := users.FindOne(context.Background(), bson.M{"name": "Alice"}).Decode(&foundUser); err != nil {
		panic(err)
	}
	println("Found user in", time.Since(start).String())

	println("Found user:", foundUser.Name, "Age:", foundUser.Age)
}
