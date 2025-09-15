package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zivwu/reminder-note-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UsersRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUsersRepository(db *mongo.Database) *UsersRepository {
	return &UsersRepository{
		db:         db,
		collection: db.Collection("users"),
	}
}

func (u *UsersRepository) UpsertUser(ctx context.Context, lineId string, name string) (*models.UserModel, error) {
	filter := bson.M{
		"lineId":  lineId,
		"deleted": false,
	}

	var result *models.UserModel

	exist := u.collection.FindOne(ctx, filter)
	if exist.Err() != nil && exist.Err() != mongo.ErrNoDocuments {
		log.Println("find user error:", exist.Err())
		return nil, exist.Err()
	}

	if exist.Err() == mongo.ErrNoDocuments {
		// 建立內容
		insertParams := models.UserModel{
			LineID:      lineId,
			Name:        name,
			Deleted:     false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			LastLoginAt: time.Now(),
		}
		insertResult, err := u.collection.InsertOne(ctx, insertParams)
		if err != nil {
			log.Println("insert user fail:", err)
			return nil, err
		}

		ID, ok := insertResult.InsertedID.(bson.ObjectID)
		if !ok {
			return nil, fmt.Errorf("parse inserted id failed")
		}

		err = u.collection.FindOne(ctx, bson.M{
			"id": ID.Hex(),
		}).Decode(&result)

		if err != nil {
			return nil, fmt.Errorf("find user after insert fail: %w ", err)
		}

		return result, nil
	} else {

		// 更新內容
		updateParams := bson.M{
			"name":        name,
			"lastLoginAt": time.Now(),
		}

		doc := bson.M{
			"$set": updateParams,
		}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err := u.collection.FindOneAndUpdate(ctx, filter, doc, opts).Decode(&result)
		if err != nil {
			log.Println("update user fail:", err)
			return nil, err
		}

		return result, nil
	}
}
