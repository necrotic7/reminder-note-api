package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zivwu/reminder-note-api/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/fx"
)

func InitMongoDB(lc fx.Lifecycle) (*mongo.Client, error) {
	connectString := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		config.Env.DB.User,
		config.Env.DB.Password,
		config.Env.DB.Host,
		config.Env.DB.Port,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(connectString))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	log.Println("MongoDB connected successfully!")

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Disconnecting MongoDB...")
			return client.Disconnect(ctx)
		},
	})
	return client, nil
}
