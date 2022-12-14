package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	dto "jasonLuFa/simpleLine-Webhook/model/DTO"
	mockdb "jasonLuFa/simpleLine-Webhook/save/mock"
	"jasonLuFa/simpleLine-Webhook/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var testAmount = 3

func TestListUserMessageAPI(t *testing.T) {
	userId := util.RandomStringAndNumber(33)
	userName := util.RandomString(5)
	testAmount := 5
	userMessages := make([]*dto.UserMessage, testAmount)
	for i := 0; i < testAmount; i++ {
		userMessages[i] = randomUserMessage(userId, userName)
	}

	testCases := []struct {
		name          string
		userId        string
		buildStubs    func(iUserMessageRepository *mockdb.MockIUserMessageRepository)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "ok",
			userId: userId,
			buildStubs: func(iUserMessageRepository *mockdb.MockIUserMessageRepository) {
				iUserMessageRepository.EXPECT().List(gomock.Any()).Times(1).Return(userMessages, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserMessages(t, recorder.Body, userMessages)
			},
		},
		{
			name:   "InternalServerError",
			userId: userId,
			buildStubs: func(iUserMessageRepository *mockdb.MockIUserMessageRepository) {
				iUserMessageRepository.EXPECT().List(gomock.Any()).Times(1).Return([]*dto.UserMessage{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			iUserMessageRepository := mockdb.NewMockIUserMessageRepository(ctrl)
			tc.buildStubs(iUserMessageRepository)

			server := newTestServer(t, iUserMessageRepository)

			url := fmt.Sprintf("/v1/users/%s/user-messages", tc.userId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestSendMsgAPI(t *testing.T) {
	userId := util.RandomStringAndNumber(33)
	sendMessages := randomMessages()

	testCases := []struct {
		name          string
		body          gin.H
		userId        string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "ok",
			body:   gin.H{"messages": sendMessages},
			userId: userId,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, {
			name:   "InvalidRequestBody",
			body:   gin.H{"messages": "invalid"},
			userId: userId,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			iUserMessageRepository := mockdb.NewMockIUserMessageRepository(ctrl)

			server := newTestServer(t, iUserMessageRepository)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest("POST", "/v1/users/:userId/user-messages/push", bytes.NewReader(data))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})

	}
}

func randomUserMessage(userId string, userName string) *dto.UserMessage {
	return &dto.UserMessage{
		Id:       util.RandomString(14),
		UserId:   userId,
		UserName: userName,
		Message:  util.RandomString(10),
	}
}

func requireBodyMatchUserMessages(t *testing.T, body *bytes.Buffer, userMessages []*dto.UserMessage) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUserMessages []*dto.UserMessage
	err = json.Unmarshal(data, &gotUserMessages)
	require.NoError(t, err)
	require.Equal(t, userMessages, gotUserMessages)
}

func randomMessages() []dto.Message {
	messages := make([]dto.Message, testAmount)
	for i := 0; i < testAmount; i++ {
		messages[i] = dto.Message{
			Type: "text",
			Text: util.RandomString(20),
		}
	}
	return messages
}
