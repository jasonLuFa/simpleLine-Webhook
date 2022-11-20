package line

import (
	api "jasonLuFa/simpleLine-Webhook/api"
	dto "jasonLuFa/simpleLine-Webhook/model/DTO"
	"jasonLuFa/simpleLine-Webhook/save"
	"jasonLuFa/simpleLine-Webhook/save/query"
	"jasonLuFa/simpleLine-Webhook/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userMessageRepository query.IUserMessageRepository
	// ctx                   context.Context
	userMessageCollection *mongo.Collection
	err                   error
	channel_access_token  string
	bot                   *linebot.Client
	config                util.Config
)

// LinebotCmd represents the linebot command
var LinebotCmd = &cobra.Command{
	Use:   "linebot",
	Short: "linebot is a palette that contains linebot based commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config, err = util.LoadConfig(".")
		if err != nil {
			log.Fatal("cannot load config", err)
		}
		mongoClient, ctx := save.InitMongo(config)

		viper.Set("channel_access_token", channel_access_token)

		bot, err = linebot.New(
			config.ChannelSecret,
			channel_access_token,
		)
		if err != nil {
			log.Fatal("new linebot err : ", err)
		}

		userMessageRepository = query.NewUserMessageRepository(mongoClient, ctx)
		server, _ := api.NewServer(userMessageRepository)
		server.Router.POST("/callback", callbackHandler())

		log.Fatal(server.Router.Run(":" + config.Port))
	},
}

func init() {
	LinebotCmd.Flags().StringVarP(&channel_access_token, "channel_access_token", "t", "", "set the channel access token")

}

// linebot callback handler
func callbackHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		events, err := bot.ParseRequest(ctx.Request)

		if err != nil {
			if err == linebot.ErrInvalidSignature {
				ctx.JSON(http.StatusBadRequest, nil)
			} else {
				ctx.JSON(http.StatusInternalServerError, nil)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				profile, err := bot.GetProfile(event.Source.UserID).Do()
				if err != nil {
					log.Print(err)
				}

				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					// Get how many remain free tier push message quota you still have this month. (maximum 500)
					quota, err := bot.GetMessageQuota().Do()
					if err != nil {
						log.Println("Quota err:", err)
						log.Println("remain Quota of free message:", quota.Value)
					}

					var userMessage dto.UserMessage
					userMessage.Id = message.ID
					userMessage.Message = message.Text
					userMessage.UserId = event.Source.UserID
					userMessage.UserName = profile.DisplayName
					userMessageRepository.Create(&userMessage)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("收到")).Do(); err != nil {
						log.Print(err)
					}
					// ... you can add more MessageType
				}
			}
		}
	}
}
