package controller

import (
	"jasonLuFa/simpleLine-Webhook/service"

	"github.com/gin-gonic/gin"
)

type UserMessageController struct {
	UserMessageService service.UserMessageService
}

func NewUserMessageController(userMessageService service.UserMessageService) UserMessageController {
	return UserMessageController{UserMessageService: userMessageService}
}

func (umc *UserMessageController) RegisterUserMessageRoutes(rg *gin.RouterGroup) {
	userMessageRoute := rg.Group("/users")
	userMessageRoute.POST("/:userId/user-messages/push", umc.UserMessageService.SendMsg)
	userMessageRoute.GET("/:userId/user-messages", umc.UserMessageService.List)
}
