package line

import (
	"context"
	"jasonLuFa/simpleLine-Webhook/controller"
	dto "jasonLuFa/simpleLine-Webhook/model/DTO"
	"jasonLuFa/simpleLine-Webhook/save/query"
	"jasonLuFa/simpleLine-Webhook/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server                *gin.Engine
	userMessageService    service.UserMessageService
	userMessageController controller.UserMessageController
	userMessageRepository query.IUserMessageRepository
	ctx                   context.Context
	userMessageCollection *mongo.Collection
	mongoClient           *mongo.Client
	err                   error
	channel_access_token  string
	bot                   *linebot.Client
	config                Config
)

// LinebotCmd represents the linebot command
var LinebotCmd = &cobra.Command{
	Use:   "linebot",
	Short: "linebot is a palette that contains linebot based commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		bot, err = linebot.New(
			config.ChannelSecret,
			channel_access_token,
		)
		viper.Set("channel_access_token", channel_access_token)

		server = gin.Default()
		userMessageRepository = query.NewUserMessageRepository(userMessageCollection, ctx)
		userMessageService = service.NewUserMessageService(userMessageRepository)
		userMessageController = controller.NewUserMessageController(userMessageService)
		server.POST("/callback", callbackHandler())

		log.Fatal(server.Run(":" + config.Port))
	},
}

func init() {
	config, err = LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	ctx = context.TODO()
	mongoConn := options.Client().ApplyURI(config.DBSource)
	mongoClient, err = mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatal("error while connecting with mongo", err)
	}

	// Know if a MongoDB server has been found and connected to
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}

	log.Println("mongo connection established")

	//
	LinebotCmd.Flags().StringVarP(&channel_access_token, "channel_access_token", "t", "", "set the channel access token")

	userMessageCollection = mongoClient.Database("userMessage").Collection("userMessages")
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

type Config struct {
	DBSource      string `mapstructure:"DB_SOURCE"`
	Port          string `mapstructure:"PORT"`
	ChannelSecret string `mapstructure:"CHANNEL_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
