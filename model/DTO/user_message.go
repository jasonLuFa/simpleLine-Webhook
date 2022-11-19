package dto

type UserMessage struct {
	Id       string `json:"id" bson:"id"`
	UserId   string `json:"userId" bson:"user_id"`
	UserName string `json:"userName" bson:"user_name"`
	Message  string `json:"message" bson:"message"`
}
