package query

import (
	"context"
	dto "jasonLuFa/simpleLine-Webhook/model/DTO"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserMessageRepository interface {
	Create(*dto.UserMessage) error
	List(string) ([]*dto.UserMessage, error)
}

type UserMessageRepository struct {
	userMessageCollection *mongo.Collection
	ctx                   context.Context
}

func NewUserMessageRepository(mongoClient *mongo.Client, ctx context.Context) IUserMessageRepository {
	userMessageCollection := mongoClient.Database("userMessage").Collection("userMessages")
	return &UserMessageRepository{
		userMessageCollection: userMessageCollection,
		ctx:                   ctx,
	}
}

func (umr *UserMessageRepository) Create(userMessage *dto.UserMessage) error {
	_, err := umr.userMessageCollection.InsertOne(umr.ctx, userMessage)
	return err
}

func (umr *UserMessageRepository) List(userId string) ([]*dto.UserMessage, error) {
	var userMessages []*dto.UserMessage
	cur, err := umr.userMessageCollection.Find(umr.ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, err
	}
	defer cur.Close(umr.ctx)

	for cur.Next(umr.ctx) {
		var userMessage dto.UserMessage
		err := cur.Decode(&userMessage)
		if err != nil {
			return nil, err
		}
		userMessages = append(userMessages, &userMessage)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	return userMessages, nil
}
