package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zivwu/reminder-note-api/internal/models"
	"github.com/zivwu/reminder-note-api/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotifyRecordsRepository struct {
	db         *mongo.Client
	collection *mongo.Collection
}

func NewNotifyRecordsRepository(db *mongo.Client) *NotifyRecordsRepository {
	return &NotifyRecordsRepository{
		db:         db,
		collection: db.Database("reminder-note").Collection("notify_records"),
	}
}

func (r *NotifyRecordsRepository) InsertNotifyRecord(ctx context.Context, params models.InsertNotifyRecord) (ID string, err error) {
	doc := models.NotifyRecordModel{
		UserID:    params.UserID,
		Content:   params.Content,
		Status:    params.Status,
		Retry:     params.Retry,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		log.Println("insert notify records fail: ", err)
		return
	}
	ID = result.InsertedID.(bson.ObjectID).Hex()
	return
}

func (r *NotifyRecordsRepository) UpdateNotifyRecord(ctx context.Context, params models.UpdateNotifyRecord) (err error) {
	// 轉換 ID 為 ObjectID
	objID, err := bson.ObjectIDFromHex(params.ID)
	if err != nil {
		log.Println("update notify records fail: ", err)
		return err
	}
	fmt.Println(params.ID, utils.ToJson(objID))
	filter := bson.M{
		"_id":    objID,
		"userId": params.UserID,
	}
	doc := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"status":    params.Status,
			"retry":     params.Retry,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, doc)
	if err != nil {
		log.Println("update notify records fail: ", err)
		return
	}
	return
}
