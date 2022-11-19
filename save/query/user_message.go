package query

import (
	"context"
	dto "jasonLuFa/simpleLine-Webhook/model/DTO"

	"go.mongodb.org/mongo-driver/mongo"
)

type IUserMessageRepository interface {
	Create(*dto.UserMessage) error
	// List(string) ([]*dto.UserMessage, error)
}

type UserMessageRepository struct {
	userMessageCollection *mongo.Collection
	ctx                   context.Context
}

func NewUserMessageRepository(userMessageCollection *mongo.Collection, ctx context.Context) IUserMessageRepository {
	return &UserMessageRepository{
		userMessageCollection: userMessageCollection,
		ctx:                   ctx,
	}
}

func (umr *UserMessageRepository) Create(userMessage *dto.UserMessage) error {
	_, err := umr.userMessageCollection.InsertOne(umr.ctx, userMessage)
	return err
}
