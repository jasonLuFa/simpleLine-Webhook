# <h1 align="center"> :bell: simpleLine-Webhook </h1>

## 💪 Purpose

- Receive message from line webhook and save the user info and message in MongoDB
- Create a API send message back to line
- Create a API query message list of the user from MongoDB

## :clapper: Demo

https://user-images.githubusercontent.com/52907691/202882712-050086ea-95f4-4380-abbf-7c0b82d6d76a.mp4

- **00:16** :point_right: input the channel access token and start the server
- **00:33 ~ 00:47** :point_right: Receive message from line webhook and save the user info and message in MongoDB
- **00:59** :point_right: Create a API query message list of the user from MongoDB
- **01:03** :point_right: Create a API send message back to line

# :toolbox: Development tool

1. postman : test API
1. ngrok : for local test to generate a https endpoint
1. mongoDB ( Studio 3T ) : mongo GUI
1. docker
1. Makefile

# :computer: Golang library

1. [Gin](https://github.com/gin-gonic/gin)
1. [Viper](https://github.com/spf13/viper)
1. [mongo-go-driver](https://github.com/mongodb/mongo-go-driver)
1. [cobra](https://github.com/spf13/cobra)

