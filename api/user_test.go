package api

import (
	mockdb "Solvery/db/mock"
	db "Solvery/db/sqlc"
	"Solvery/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateUserApi(t *testing.T) {
	user := db.User{
		Name:      util.RandomName(),
		Class:     util.RandomClass(),
		Email:     util.RandomEmail(),
		Credit:    0,
		CreatedAt: time.Now(),
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"full_name": user.Name,
				"class":     user.Class,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Name:  user.Name,
					Class: user.Class,
					Email: user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), arg).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"full_name": user.Name,
				"class":     user.Class,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Name:  user.Name,
					Class: user.Class,
					Email: user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), arg).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateEmail",
			body: gin.H{
				"full_name": user.Name,
				"class":     user.Class,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Name:  user.Name,
					Class: user.Class,
					Email: user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), arg).Times(1).Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"full_name": "",
				"class":     user.Class,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidClass",
			body: gin.H{
				"full_name": user.Name,
				"class":     "",
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"full_name": user.Name,
				"class":     user.Class,
				"email":     "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
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

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetUserApi(t *testing.T) {
	user := db.User{
		Name:      util.RandomName(),
		Class:     util.RandomClass(),
		Email:     util.RandomEmail(),
		Credit:    0,
		CreatedAt: time.Now(),
	}

	testCases := []struct {
		name          string
		email         string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			email: user.Email,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), user.Email).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:  "InternalError",
			email: user.Email,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), user.Email).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "NotFound",
			email: "somenotexisting@example.com",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "InvalidEmail",
			email: "invalid-email",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
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

			url := fmt.Sprintf("/users/%s", tc.email)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListUsersApi(t *testing.T) {
	testCases := []struct {
		name          string
		pageId        int
		pageSize      int
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageId:   1,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListUsersParams{
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().ListUsers(gomock.Any(), arg).Times(1).Return([]db.User{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			pageId:   1,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListUsersParams{
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().ListUsers(gomock.Any(), arg).Times(1).Return([]db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidPageId",
			pageId:   0,
			pageSize: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "InvalidPageSize",
			pageId:   1,
			pageSize: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Times(0)
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

			url := fmt.Sprintf("/users?page_id=%d&page_size=%d", tc.pageId, tc.pageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateUserApi(t *testing.T) {
	user := db.User{
		Name:      util.RandomName(),
		Class:     util.RandomClass(),
		Email:     util.RandomEmail(),
		Credit:    0,
		CreatedAt: time.Now(),
	}

	amount := int32(10)

	comment := "update credit"

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"amount": amount,
				"email":  user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.PaymentTxParams{
					Amount:    amount,
					UserEmail: user.Email,
					Comment:   comment,
				}
				store.EXPECT().GetUser(gomock.Any(), user.Email).Times(1).Return(user, nil)
				store.EXPECT().PaymentTx(gomock.Any(), arg).Times(1).Return(db.PaymentTxResult{
					User: db.User{
						Name:      user.Name,
						Class:     user.Class,
						Email:     user.Email,
						Credit:    user.Credit + amount,
						CreatedAt: user.CreatedAt,
					},
					Entry: db.Entry{
						Amount:    amount,
						Comment:   comment,
						UserEmail: user.Email,
					},
				}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPaymentTxResult(t, recorder.Body, db.PaymentTxResult{
					User: db.User{
						Name:      user.Name,
						Class:     user.Class,
						Email:     user.Email,
						Credit:    user.Credit + amount,
						CreatedAt: user.CreatedAt,
					},
					Entry: db.Entry{
						Amount:    amount,
						Comment:   comment,
						UserEmail: user.Email,
					},
				})
			},
		},
		{
			name: "InternalError1",
			body: gin.H{
				"amount": amount,
				"email":  user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.PaymentTxParams{
					Amount:    amount,
					UserEmail: user.Email,
					Comment:   comment,
				}
				store.EXPECT().GetUser(gomock.Any(), user.Email).Times(1).Return(db.User{}, sql.ErrConnDone)
				store.EXPECT().PaymentTx(gomock.Any(), arg).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalError2",
			body: gin.H{
				"amount": amount,
				"email":  user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.PaymentTxParams{
					Amount:    amount,
					UserEmail: user.Email,
					Comment:   comment,
				}
				store.EXPECT().GetUser(gomock.Any(), user.Email).Times(1).Return(user, nil)
				store.EXPECT().PaymentTx(gomock.Any(), arg).Times(1).Return(db.PaymentTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NotFound",
			body: gin.H{
				"amount": amount,
				"email":  user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), user.Email).Times(1).Return(db.User{}, sql.ErrNoRows)
				store.EXPECT().PaymentTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"amount": amount,
				"email":  "invalid-email",
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
				"amount": -1000000,
				"email":  user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), user.Email).Times(1).Return(user, nil)
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

			url := "/users/update"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			//require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Name, gotUser.Name)
	require.Equal(t, user.Class, gotUser.Class)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.Credit, gotUser.Credit)
	require.WithinDuration(t, user.CreatedAt, gotUser.CreatedAt, time.Second)

}

func requireBodyMatchPaymentTxResult(t *testing.T, body *bytes.Buffer, result db.PaymentTxResult) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotResult db.PaymentTxResult
	err = json.Unmarshal(data, &gotResult)

	require.NoError(t, err)
	require.Equal(t, result.User.Name, gotResult.User.Name)
	require.Equal(t, result.User.Class, gotResult.User.Class)
	require.Equal(t, result.User.Email, gotResult.User.Email)
	require.Equal(t, result.User.Credit, gotResult.User.Credit)
	require.WithinDuration(t, result.User.CreatedAt, gotResult.User.CreatedAt, time.Second)

	require.Equal(t, result.Entry.Amount, gotResult.Entry.Amount)
	require.Equal(t, result.Entry.Comment, gotResult.Entry.Comment)
	require.Equal(t, result.Entry.UserEmail, gotResult.Entry.UserEmail)
}
