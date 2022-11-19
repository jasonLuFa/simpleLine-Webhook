package controller

import (
	"jasonLuFa/simpleLine-Webhook/service"
)

type UserMessageController struct {
	UserMessageService service.UserMessageService
}

func NewUserMessageController(userMessageService service.UserMessageService) UserMessageController {
	return UserMessageController{UserMessageService: userMessageService}
}
