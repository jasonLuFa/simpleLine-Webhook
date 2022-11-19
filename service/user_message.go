package service

import (
	"jasonLuFa/simpleLine-Webhook/save/query"
)

type UserMessageService struct {
	IUserMessageRepository query.IUserMessageRepository
}

func NewUserMessageService(iUserMessageRepository query.IUserMessageRepository) UserMessageService {
	return UserMessageService{
		IUserMessageRepository: iUserMessageRepository,
	}
}
