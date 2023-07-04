package api

import (
	mockdb "Solvery/db/mock"
	db "Solvery/db/sqlc"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTas1Api(t *testing.T) {
	array := []int{1, 2, 2, 5, 6, 5, 6}
	email := "test@example.com"

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"array": array,
				"email": email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.PaymentTxParams{
					Amount:    -int32(len(array)),
					UserEmail: email,
					Comment:   fmt.Sprintf("%s, input: %v", task1Comment, array),
				}
				store.EXPECT().GetUser(gomock.Any(), email).Times(1).Return(db.User{Email: email}, nil)
				store.EXPECT().PaymentTx(gomock.Any(), arg).Times(1).Return(db.PaymentTxResult{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError1",
			body: gin.H{
				"array": array,
				"email": email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), email).Times(1).Return(db.User{}, sql.ErrConnDone)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalError2",
			body: gin.H{
				"array": array,
				"email": email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.PaymentTxParams{
					Amount:    -int32(len(array)),
					UserEmail: email,
					Comment:   fmt.Sprintf("%s, input: %v", task1Comment, array),
				}
				store.EXPECT().GetUser(gomock.Any(), email).Times(1).Return(db.User{Email: email}, nil)
				store.EXPECT().PaymentTx(gomock.Any(), arg).Times(1).Return(db.PaymentTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"array": array,
				"email": email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), email).Times(1).Return(db.User{}, sql.ErrNoRows)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"array": array,
				"email": "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidArray1",
			body: gin.H{
				"array": []int{1, 2, 2, 5, 6, 5, 6, 100},
				"email": "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidArray2",
			body: gin.H{
				"array": []string{"1", "2", "2", "5", "6", "5", "6"},
				"email": "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "lowCredit",
			body: gin.H{
				"array": []int{1, 2, 2, 5, 6, 5, 6},
				"email": email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), email).Times(1).Return(db.User{Credit: -994}, nil)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
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

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/task1"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			//require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
