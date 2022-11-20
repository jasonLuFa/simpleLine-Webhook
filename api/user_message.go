package api

import (
	"bytes"
	"encoding/json"
	dto "jasonLuFa/simpleLine-Webhook/model/DTO"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (server *Server) SendMsg(c *gin.Context) {
	var sendMsg dto.SendMessage
	if err := c.ShouldBindJSON(&sendMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		log.Println(err.Error())
		return
	}

	userId := c.Param("userId")
	res, err := server.sendMessageRequest(userId, sendMsg)
	if err != nil {
		c.JSON(res.StatusCode, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/json")
	defer res.Body.Close()
	c.JSON(200, res.Body)
}

func (server *Server) sendMessageRequest(userId string, sendMsg dto.SendMessage) (*http.Response, error) {
	sendMsg.TO = userId
	sendMsgJson, err := json.Marshal(sendMsg)
	if err != nil {
		log.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}
	log.Printf("sendMessage JSON: %s\n", string(sendMsgJson))
	body := []byte(sendMsgJson)

	r, err := http.NewRequest(http.MethodPost, "https://api.line.me/v2/bot/message/push", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error occured during make a request of line API. Error: %s", err.Error())
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+viper.GetViper().GetString("channel_access_token"))
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Fatalf("Error occured during get the resposne of the line API. Error: %s", err.Error())
	}

	return res, err
}

func (server *Server) List(ctx *gin.Context) {
	userId := ctx.Param("userId")
	users, err := server.iUserMessageRepository.List(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
