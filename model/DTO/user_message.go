package dto

type UserMessage struct {
	Id       string `json:"id" bson:"id"`
	UserId   string `json:"userId" bson:"user_id"`
	UserName string `json:"userName" bson:"user_name"`
	Message  string `json:"message" bson:"message"`
}

type SendMessage struct {
	TO       string    `json:"to" bson:"to"`
	Messages []Message `json:"messages" bson:"messages"`
}

type Message struct {
	Type string `json:"type" bson:"type"`
	Text string `json:"text" bson:"text"`
}
