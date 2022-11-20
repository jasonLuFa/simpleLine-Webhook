package api

import (
	"jasonLuFa/simpleLine-Webhook/save/query"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router                 *gin.Engine
	iUserMessageRepository query.IUserMessageRepository
}

func NewServer(iUserMessageRepository query.IUserMessageRepository) (*Server, error) {

	server := &Server{iUserMessageRepository: iUserMessageRepository}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	basePath := router.Group("/v1")

	userMessageRoute := basePath.Group("/users")
	userMessageRoute.POST("/:userId/user-messages/push", server.SendMsg)
	userMessageRoute.GET("/:userId/user-messages", server.List)

	server.Router = router
}
