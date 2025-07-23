package repositories

import (
	"context"
	"log"
	"time"

	"github.com/zivwu/reminder-note-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotifyRecordsRepository struct {
	db *mongo.Client
}

func NewNotifyRecordsRepository(db *mongo.Client) *NotifyRecordsRepository {
	return &NotifyRecordsRepository{
		db: db,
	}
}

func (r *NotifyRecordsRepository) InsertNotifyRecord(ctx context.Context, params models.InsertNotifyRecord) (ID string, err error) {
	collection := r.db.Database("reminder-note").Collection("notify_records")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	doc := models.NotifyRecordModel{
		UserID:  params.UserID,
		Content: params.Content,
		Status:  params.Status,
		Retry:   params.Retry,
	}

	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Println("insert reminder fail: ", err)
		return
	}
	ID = result.InsertedID.(bson.ObjectID).String()
	return
}
