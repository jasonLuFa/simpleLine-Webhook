package service

import (
	"bytes"
	"encoding/json"
	dto "jasonLuFa/simpleLine-Webhook/model/DTO"
	"jasonLuFa/simpleLine-Webhook/save/query"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type UserMessageService struct {
	IUserMessageRepository query.IUserMessageRepository
}

func NewUserMessageService(iUserMessageRepository query.IUserMessageRepository) UserMessageService {
	return UserMessageService{
		IUserMessageRepository: iUserMessageRepository,
	}
}

func (ums *UserMessageService) SendMsg(c *gin.Context) {
	var sendMsg dto.SendMessage
	if err := c.ShouldBindJSON(&sendMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		log.Println(err.Error())
		return
	}

	userId := c.Param("userId")
	res, err := ums.sendMessageRequest(userId, sendMsg)
	if err != nil {
		c.JSON(res.StatusCode, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/json")
	defer res.Body.Close()
	c.JSON(200, res.Body)
}

func (ums *UserMessageService) sendMessageRequest(userId string, sendMsg dto.SendMessage) (*http.Response, error) {
	sendMsg.TO = userId
	sendMsgJson, err := json.Marshal(sendMsg)
	if err != nil {
		log.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}
	log.Printf("sendMessage JSON: %s\n", string(sendMsgJson))
	body := []byte(sendMsgJson)

	r, err := http.NewRequest(http.MethodPost, "https://api.line.me/v2/bot/message/push", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+viper.GetViper().GetString("channel_access_token"))
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	return res, err
}

func (ums *UserMessageService) List(ctx *gin.Context) {
	userId := ctx.Param("userId")
	users, err := ums.IUserMessageRepository.List(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
