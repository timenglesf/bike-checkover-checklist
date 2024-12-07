package db

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(uri string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(
		options.ServerAPIVersion1,
	)

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	var result bson.M

	if err := client.Database("admin").RunCommand(
		context.TODO(),
		bson.D{{
			"ping",
			1,
		}}).
		Decode(&result); err != nil {
		return nil, err
	}

	return client, nil
}

func Close(client *mongo.Client) error {
	return client.Disconnect(context.Background())
}

func ConnectWithRetries(logger *slog.Logger, attempts int) (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		logger.Error("MONGO_URI is not set")
		return nil, ErrDBURIEmpty
	}

	var err error
	var client *mongo.Client
	for i := 0; i < attempts; i++ {
		client, err = Connect(mongoURI)
		if err == nil {
			break
		}
		logger.Info("Retrying database connection", "attempt", i+1, "err", err.Error())
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		logger.Error("Failed to connect to database", "err", err.Error())
		return nil, err
	}
	logger.Info("Database connection successful")
	return client, nil
}
