mongo:
	docker run -d --name mongo4.4 -p 27017:27017  mongo:4.4

mock:
	mockgen -package mockdb -destination save/mock/user_message.go jasonLuFa/simpleLine-Webhook/save/query IUserMessageRepository

test:
	go test -v -cover ./...

.PHONY: mongo mock test