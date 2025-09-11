package repositories

import (
	"context"

	"github.com/zivwu/reminder-note-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CountersRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewCountersRepository(db *mongo.Database) *CountersRepository {
	return &CountersRepository{
		db:         db,
		collection: db.Collection("counters"),
	}
}

// TODO 暫未啟用
func (c *CountersRepository) GetNextSeq(ctx context.Context, model models.EnumCounterModel) (*int, error) {
	var result models.CounterModel
	err := c.collection.FindOneAndUpdate(
		ctx,
		bson.M{"model": model},
		bson.M{"$inc": bson.M{"seq": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result.Seq, nil
}
