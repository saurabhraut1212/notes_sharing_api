package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //create context with 10 sec, Ensures that if MongoDB doesn’t respond within 10 seconds, the attempt is canceled.
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil { //Sends a ping to MongoDB to confirm the connection is active.
		return nil, err
	}

	return client, nil
}
